package btcplex

import (
    "log"
    "fmt"
    _ "github.com/jmhodges/levigo"
    "net/http"
    "strconv"
    "strings"
    "time"
    "io"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/codegangsta/martini-contrib/binding"
    "github.com/pmylund/go-cache"
    "html/template"
    "github.com/garyburd/redigo/redis"
)

// Martini form for the search input
type SearchForm struct {
    Query string `form:"q"`
}

// Struct holding page meta data, like meta tags, and some template variables
type PageMeta struct {
    Title string
    Description string
    Menu string
    Block *Block
    Blocks *[]*Block
    Tx *Tx
    AddressData *AddressData
    LastHeight uint
    CurrentHeight uint
    Error string
}

type RedisWrapper struct {
    Pool *redis.Pool
}

const (
    ratelimitwindow = 3600
    ratelimitcnt = 3600
)

func RateLimited(rediswrapper *RedisWrapper, ip string) (bool, int, int) {
    conn := rediswrapper.Pool.Get()
    defer conn.Close()
    reset := int(time.Now().UTC().Unix() / ratelimitwindow * ratelimitwindow + ratelimitwindow)
    ipkey := fmt.Sprintf("rl:%v:%v", ip, reset)
    cnt, _ := redis.Int(conn.Do("GET", ipkey))
    if cnt > ratelimitcnt {
        return true, cnt, reset
    } else {
        conn.Send("MULTI")
        conn.Send("INCR", ipkey)
        conn.Send("EXPIREAT", ipkey, reset + ratelimitwindow)
        conn.Do("EXEC")
        cnt+=1
        return false, cnt, reset
    }
}

func Run() {
    log.Println("Starting")
    conf, _ := LoadConfig("./config.json")

    // Setting up cache
    c := cache.New(15*time.Minute, 30*time.Second)

    // Redis connect
    // Used for pub/sub in the webapp and data like latest processed height
    pool, _ := GetRedis(conf)
    rediswrapper := new(RedisWrapper)
    rediswrapper.Pool = pool
    ssdb, _ := GetSSDB(conf)

    //blocknotifylive := 0
    //ticker := time.NewTicker(time.Second * 3)
    //go func() {
    //    for _ = range ticker.C {
    //        fmt.Printf("blocknotify live: %v\n", blocknotifylive)
    //    }
    //}()

    AppHelpers := template.FuncMap{
        "cut": func(addr string, length int) string {
            return fmt.Sprintf("%v...", addr[:length])
        },
        "cutmiddle": func(addr string, length int) string {
            return fmt.Sprintf("%v...%v", addr[:length], addr[len(addr) - length:])
        },
        "tokb": func(size uint32) string {
            return fmt.Sprintf("%.3f", float32(size) / 1024)
        },
        "computefee": func(tx *Tx) string {
            if tx.TotalIn == 0 {
                return "0"
            }
            return fmt.Sprintf("%v", float32(tx.TotalIn - tx.TotalOut) / 1e8)
        },
        "generationmsg": func(tx *Tx) string {
            reward := GetBlockReward(tx.BlockHeight)
            fee := float64(tx.TotalOut - uint64(reward)) / 1e8
            return fmt.Sprintf("%v BTC + %v total fees", float64(reward) / 1e8, fee)
        },
        "tobtc": func(val uint64) string {
            return fmt.Sprintf("%v", float64(val) / 1e8)
        },
        "inttobtc": func(val int64) string {
            return fmt.Sprintf("%v", float64(val) / 1e8)
        },
        "formatprevout": func(prevout *PrevOut) string {
            return fmt.Sprintf("%v:%v", prevout.Hash, prevout.Vout)
        },
        "formattime": func(ts uint32) string {
            return fmt.Sprintf("%v", time.Unix(int64(ts), 0).UTC())
        },
        "formatiso": func(ts uint32) string {
            return fmt.Sprintf("%v", time.Unix(int64(ts), 0).Format(time.RFC3339))
        },
        "sub": func(h, p uint) uint {
            return h - p
        },
        "add": func(h, p uint) uint {
            return h + p
        },
        "confirmation": func(height uint) uint {
            //lh, _ := c.Get("last-height")
            lh := 200000
            //return lh - height + 1
            return uint(lh)
        },

    }

    conn := ssdb.Get()    
    latestheight, _ := redis.Int(conn.Do("GET", "height:latest"))
    log.Printf("Latest height: %v\n", latestheight)
    //lblock, _ := GetBlockByHeight(db, uint(latestheight))
    latesthash := ""
    conn.Close()

    m := martini.Classic()
    m.Map(c)
    m.Map(rediswrapper)
    m.Map(ssdb)
    m.Use(render.Renderer(render.Options{
        Layout: "layout",
        Funcs: []template.FuncMap{AppHelpers},
    }))

    // We rate limit the API
    m.Use(func(res http.ResponseWriter, req *http.Request, rediswrapper *RedisWrapper, log *log.Logger) {
        remoteIP := strings.Split(req.RemoteAddr,":")[0]
        _, xforwardedfor := req.Header["X-Forwarded-For"]
        if xforwardedfor {
            remoteIP = req.Header["X-Forwarded-For"][0]
        }
        log.Printf("R:%v\nip:%+v\n", time.Now(), remoteIP)
        if strings.Contains(req.RequestURI, "/api/v") {
            ratelimited, cnt, reset := RateLimited(rediswrapper, remoteIP)
            res.Header().Set("X-RateLimit-Limit", strconv.Itoa(ratelimitcnt))
            res.Header().Set("X-RateLimit-Remaining", strconv.Itoa(ratelimitcnt - cnt))
            res.Header().Set("X-RateLimit-Reset", strconv.Itoa(reset))
            // Set CORS header
            res.Header().Set("Access-Control-Expose-Headers", " X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
            res.Header().Set("Access-Control-Allow-Origin", "*")
            if ratelimited {
                res.WriteHeader(429)
            }
        }
    })

    m.Get("/robots.txt", func() string {
        return "User-agent: *\nDisallow: /api/v1"
    })

    m.Get("/", func(r render.Render, c *cache.Cache, db *redis.Pool) {
        pm := new(PageMeta)
        blocks, _ := GetLastXBlocks(db, uint(latestheight), uint(latestheight - 30))
        pm.Blocks = &blocks
        pm.Title = "Latest Bitcoin blocks"
        pm.Description = "Open source Bitcoin block chain explorer with JSON API"
        pm.Menu = "latest_blocks"
        pm.LastHeight = uint(latestheight)
        r.HTML(200, "index", &pm)
    })

    m.Get("/blocks/:currentheight", func(params martini.Params, r render.Render, db *redis.Pool) {
        pm := new(PageMeta)
        currentheight, _ := strconv.ParseUint(params["currentheight"], 10, 0)
        blocks, _ := GetLastXBlocks(db, uint(currentheight), uint(currentheight - 30))
        pm.Blocks = &blocks
        pm.Title = "Bitcoin blocks"
        pm.Menu = "blocks"
        pm.LastHeight = uint(latestheight)
        pm.CurrentHeight = uint(currentheight)
        r.HTML(200, "blocks", &pm)
    })

    m.Get("/block/:hash", func(params martini.Params, r render.Render, db *redis.Pool) {
        pm := new(PageMeta)
        block, _ := GetBlockByHash(db, params["hash"])
        block.FetchTxs(db)
        pm.Block = block
        pm.Title = fmt.Sprintf("Bitcoin block #%v", block.Height)
        pm.Description = fmt.Sprintf("Bitcoin block #%v summary and related transactions", block.Height)
        r.HTML(200, "block", &pm)
    })
    m.Get("/api/v1/block/:hash", func(params martini.Params, r render.Render, db *redis.Pool) {
        block, _ := GetBlockByHash(db, params["hash"])
        block.FetchTxs(db)
        r.JSON(200, block)
    })

    m.Get("/tx/:hash", func(params martini.Params, r render.Render, db *redis.Pool) {
        pm := new(PageMeta)
        tx, _ := GetTx(db, params["hash"])
        tx.Build(db)
        pm.Tx = tx
        pm.Title = fmt.Sprintf("Bitcoin transaction %v", tx.Hash)
        pm.Description = fmt.Sprintf("Bitcoin transaction %v summary.", tx.Hash)
        r.HTML(200, "tx", pm)
    })
    m.Get("/api/v1/tx/:hash", func(params martini.Params, r render.Render, db *redis.Pool) {
        tx, _ := GetTx(db, params["hash"])
        tx.Build(db)
        r.JSON(200, tx)
    })

    m.Get("/address/:address", func(params martini.Params, r render.Render, db *redis.Pool) {
        pm := new(PageMeta)
        addressdata, _ := GetAddress(db, params["address"])
        pm.AddressData = addressdata
        pm.Title = fmt.Sprintf("Bitcoin address %v", params["address"])
        pm.Description = fmt.Sprintf("Transactions and summary for the Bitcoin address %v.", params["address"])
        r.HTML(200, "address", pm)
    })
    m.Get("/api/v1/address/:address", func(params martini.Params, r render.Render, db *redis.Pool) {
        addressdata, _ := GetAddress(db, params["address"])
        r.JSON(200, addressdata)
    })

    m.Get("/api", func(r render.Render) {
        pm := new(PageMeta)
        pm.Title = "API Documentation"
        pm.Description = "BTCPlex provides JSON API for developers to retrieve Bitcoin block chain data pragmatically"
        pm.Menu = "api"
        r.HTML(200, "api_docs", pm)
    })

    m.Get("/about", func(r render.Render) {
        pm := new(PageMeta)
        pm.Title = "About"
        pm.Description = "Learn more about BTCPlex, an open source Bitcoin block chain explorer with JSON API"
        pm.Menu = "about"
        r.HTML(200, "about", pm)
    })

    m.Post("/search", binding.Form(SearchForm{}), binding.ErrorHandler, func(search SearchForm, r render.Render, db *redis.Pool) {
        pm := new(PageMeta)
        // Check if the query isa block height
        isblockheight, hash := IsBlockHeight(db, search.Query)
        if isblockheight && hash != "" {
            r.Redirect(fmt.Sprintf("/block/%v", hash))
        }
        // Check if the query is block hash
        isblockhash, hash := IsBlockHash(db, search.Query)
        if isblockhash {
            r.Redirect(fmt.Sprintf("/block/%v", hash))    
        }
        // Check for TX
        istxhash, txhash := IsTxHash(db, search.Query)
        if istxhash {
            r.Redirect(fmt.Sprintf("/tx/%v", txhash))
        }
        // Check for Bitcoin address
        isaddress, address := IsAddress(search.Query)
        if isaddress {
            r.Redirect(fmt.Sprintf("/address/%v", address))
        }
        pm.Title = "Search"
        pm.Error = "Nothing found"
        r.HTML(200, "search", pm)
    })

    m.Get("/api/v1/getblockcount", func(r render.Render) {
        r.JSON(200, latestheight)
    })

    m.Get("/api/v1/latesthash", func(r render.Render) {
        r.JSON(200, latesthash)
    })

    m.Get("/api/v1/getblockhash/:height", func(r render.Render, params martini.Params, db *redis.Pool) {
        height, _ := strconv.ParseUint(params["height"], 10, 0)
        blockhash, _ := GetBlockHash(db, uint(height))
        r.JSON(200, blockhash)
    })

    m.Get("/api/v1/getreceivedbyaddress/:address", func(r render.Render, params martini.Params, db *redis.Pool) {
        res, _ := GetReceivedByAddress(db, params["address"])
        r.JSON(200, res)
    })

    m.Get("/api/v1/getsentbyaddress/:address", func(r render.Render, params martini.Params, db *redis.Pool) {
        res, _ := GetSentByAddress(db, params["address"])
        r.JSON(200, res)
    })

    m.Get("/api/v1/addressbalance/:address", func(r render.Render, params martini.Params, db *redis.Pool) {
        res, _ := AddressBalance(db, params["address"])
        r.JSON(200, res)
    })

    m.Get("/api/v1/checkaddress/:address", func(params martini.Params, r render.Render) {
        valid, _ := ValidA58([]byte(params["address"]))
        r.JSON(200, valid)
    })

    m.Get("/api/v1/blocknotify", func(w http.ResponseWriter, r *http.Request, pool *RedisWrapper) {
        //blocknotifylive += 1
        // TODO remplacer ca par un INCR
        conn := rediswrapper.Pool.Get()
        defer conn.Close()
        //defer func() {
        //    blocknotifylive -= 1
        //}()
        psc := redis.PubSubConn{Conn: conn}
        running := true
        notifier := w.(http.CloseNotifier).CloseNotify()
        timer := time.NewTimer(time.Second *1300)

        f, _ := w.(http.Flusher)
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

        psc.Subscribe("btcplex:blocknotify")
        rec := make(chan string)
        go func(rec chan string) {
            for {
                switch v := psc.Receive().(type) {
                case redis.Message:
                    rec <-string(v.Data)
                }
            }
        }(rec)
        var ls string
        for {
            if running {
                select {
                    case ls = <-rec:
                        io.WriteString(w, fmt.Sprintf("data: %v\n\n", ls))
                        f.Flush()
                    case <-notifier:
                        running = false
                        log.Println("CLOSED")
                        break
                    case <-timer.C:
                        running = false
                        log.Println("TimeOUT")
                }
            } else {
                log.Println("DONE")
                break
            }
        }
    })

    http.ListenAndServe(fmt.Sprintf(":%v", conf.AppPort), m)
}
