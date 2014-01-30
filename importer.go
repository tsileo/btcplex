package main
import (
    "log"
    "blkparser"
    "sync"
    "time"
    "fmt"
    "github.com/garyburd/redigo/redis"
    "os"
    "os/signal"
    "encoding/json"
    "github.com/pmylund/go-cache"
    "github.com/jmhodges/levigo"
    "runtime"
    _ "btcplex"
)

type Block struct {
    Hash        string  `json:"hash"`
    Height      uint    `json:"height"`
//    Txs         []*Tx   `json:"tx,omitempty" bson:"-"`
    Version     uint32  `json:"ver"`
    MerkleRoot  string  `json:"mrkl_root"`
    BlockTime   uint32  `json:"time"`
    Bits        uint32  `json:"bits"`
    Nonce       uint32  `json:"nonce"`
    Size        uint32  `json:"size"`
    TxCnt       uint32  `json:"n_tx"`
    TotalBTC    uint64  `json:"total_out"`
//    BlockReward float64 `json:"-"`
    Parent      string  `json:"prev_block"`
//    Next        string  `json:"next_block"`
}

type Tx struct {
    Hash          string         `json:"hash"`
    Index uint32 `json:"-"`
    Size          uint32         `json:"size"`
    LockTime      uint32         `json:"lock_time"`
    Version       uint32         `json:"ver"`
    TxInCnt       uint32         `json:"vin_sz"`
    TxOutCnt      uint32         `json:"vout_sz"`
    TxIns         []*TxIn        `json:"in" bson:"-"`
    TxOuts        []*TxOut       `json:"out" bson:"-"`
    TotalOut      uint64         `json:"vout_total"`
    TotalIn       uint64         `json:"vin_total"`
    BlockHash     string         `json:"block_hash"`
    BlockHeight   uint           `json:"block_height"`
    BlockTime     uint32         `json:"block_time"`
}

type TxOut struct {
    TxHash     string    `json:"-"`
    BlockHash     string    `json:"-"`
    BlockTime     uint32    `json:"-"`
    Addr     string    `json:"hash"`
    Value    uint64    `json:"value"`
    Index    uint32    `json:"n"`
    Spent    *TxoSpent `json:"spent,omitempty"`
}

type TxOutCached struct {
    Id string `json:"id"`
    Addr     string    `json:"hash"`
    Value    uint64    `json:"value"`
}

type PrevOut struct {
    Hash    string `json:"hash"`
    Vout    uint32 `json:"n"`
    Address string `json:"address"`
    Value   uint64 `json:"value"`
}


type TxIn struct {
    TxHash     string    `json:"-"`
    BlockHash     string    `json:"-"`
    BlockTime     uint32    `json:"-"`
    PrevOut   *PrevOut `json:"prev_out"`   
    Index    uint32    `json:"n"`
}

type TxoSpent struct {
    Spent       bool   `json:"spent"`
    BlockHeight uint32 `json:"block_height,omitempty" gorethink:",omitempty"`
    InputHash   string `json:"tx_hash,omitempty"  gorethink:",omitempty"`
    InputIndex  uint32 `json:"in_index,omitempty" gorethink:",omitempty"`
}

var wg sync.WaitGroup
var running bool

func getGOMAXPROCS() int {
    return runtime.GOMAXPROCS(0)
}

func main () {
    fmt.Printf("GOMAXPROCS is %d\n", getGOMAXPROCS())
    opts := levigo.NewOptions()
    opts.SetCreateIfMissing(true)
    filter := levigo.NewBloomFilter(10)
    opts.SetFilterPolicy(filter)
    ldb, err := levigo.Open("/home/thomas/btcplex_txocached100", opts) //alpha

    defer ldb.Close()

    if err != nil {
        fmt.Printf("failed to load db: %s\n", err)
    }
    
    wo := levigo.NewWriteOptions()
    //wo.SetSync(true)
    defer wo.Close()

    ro := levigo.NewReadOptions()
    defer ro.Close()

    wb := levigo.NewWriteBatch()
    defer wb.Close()

    // Redis connect
    // Used for pub/sub in the webapp and data like latest processed height
    server := "localhost:6381"
    pool := &redis.Pool{
             MaxIdle: 3,
             IdleTimeout: 240 * time.Second,
             Dial: func () (redis.Conn, error) {
                 c, err := redis.Dial("tcp", server)
                 if err != nil {
                     return nil, err
                 }
//                 if _, err := c.Do("AUTH", password); err != nil {
//                     c.Close()
//                     return nil, err
//               }
//                 return c, err
                return c, err
             },
                TestOnBorrow: func(c redis.Conn, t time.Time) error {
                    _, err := c.Do("PING")
                 return err
                },
         }
    conn := pool.Get()
    defer conn.Close()

    //latestheight, _ := redis.Int(conn.Do("GET", "height:latest"))
    latestheight := 0
    log.Printf("Latest height: %v\n", latestheight)

    running = true
    cs := make(chan os.Signal, 1)
    signal.Notify(cs, os.Interrupt)
    go func() {
        for sig := range cs {
            running = false
            log.Printf("Captured %v, waiting for everything to finish...\n", sig)
            wg.Wait()
            defer os.Exit(1)
        }
    }()

    c := cache.New(5*time.Minute, 30*time.Second)

    log.Println("DB Loaded")


    //concurrency := 50
    //sem := make(chan bool, concurrency)
    //for i := 0; i < cap(sem); i++ {
    //    sem <- true
    //}

    // Real network magic byte
    blockchain, _ := blkparser.NewBlockchain("/box/bitcoind_data/blocks", [4]byte{0xF9,0xBE,0xB4,0xD9})

    block_height := uint(0)
    //if latestheight != 0 {
        //65824908, 42014/01/29 13:05:13 Current block: 145583
    //:33764696, 12014/01/30 01:18:34 Current block: 124609
    //    err = blockchain.SkipTo(uint32(1), int64(33764696))
    //    block_height = 124609
    //    if err != nil {
    //        log.Println("Error blkparser: blockchain.SkipTo")
    //        os.Exit(1)
    //    }
    //}
    for i := uint(0); i < 280000; i++ {
        if !running {
            break
        }

        bl, er := blockchain.NextBlock()
        if er!=nil {
            wg.Wait()
            log.Println("END of DB file")
            break
        }

        bl.Raw = nil

        if bl.Parent == "" {
            block_height = uint(0)
        } else {
            prev_height, found := c.Get(bl.Parent)
            if found {
                block_height = uint(prev_height.(uint) + 1)
            }
        }
        c.Set(bl.Hash, block_height, 10*time.Minute)

        if latestheight != 0 && !(latestheight + 1 <= int(block_height)) {
            log.Printf("Skipping block #%v\n", block_height)
            continue
        }
        
        wg.Add(1)

        log.Printf("Current block: %v (%v)\n", block_height, bl.Hash)
        
        block := new(Block)
        block.Hash = bl.Hash
        block.Height = block_height
        block.Version = bl.Version
        block.MerkleRoot = bl.MerkleRoot
        block.BlockTime = bl.BlockTime
        block.Bits = bl.Bits
        block.Nonce = bl.Nonce
        block.Size = bl.Size
        block.Parent = bl.Parent

        txs := []*Tx{}

        total_bl_out := uint(0)
        for tx_index, tx := range bl.Txs {
            log.Printf("Tx #%v: %v\n", tx_index, tx.Hash)
            
            total_tx_out := uint(0)
            total_tx_in := uint(0)

            log.Println("Starting TXOS")


            //conn.Send("MULTI")
            txos := []*TxOut{}
            for txo_index, txo := range tx.TxOuts {
                total_tx_out+= uint(txo.Value)

                ntxo := new(TxOut)
                ntxo.TxHash = tx.Hash
                ntxo.BlockHash = bl.Hash
                ntxo.BlockTime = bl.BlockTime
                ntxo.Addr = txo.Addr
                ntxo.Value = txo.Value
                ntxo.Index = uint32(txo_index)
                txospent := new(TxoSpent)
                ntxo.Spent = txospent
                txos = append(txos, ntxo)

                ntxocached := new(TxOutCached)
                ntxocached.Addr = txo.Addr
                ntxocached.Value = txo.Value

                ntxocachedjson, _ := json.Marshal(ntxocached)
                wb.Put([]byte(fmt.Sprintf("txo:%v:%v", tx.Hash, txo_index)), ntxocachedjson)

                //log.Println("Before insert txo")
                //wg.Add(1)
                //sem <-true
                //go func(ntxo *TxOut) {
                //    defer wg.Done()
                //    defer func() { <-sem }()
                //    db.C("txos").Insert(ntxo)
                //}(ntxo)
                //log.Println("After insert txo")
                ntxojson, _ := json.Marshal(ntxo)
                ntxokey := fmt.Sprintf("txo:%v:%v", tx.Hash, txo_index)
                conn.Do("SET", ntxokey, ntxojson)
                //conn.Send("ZADD", fmt.Sprintf("txo:%v", tx.Hash), txo_index, ntxokey)
                conn.Do("ZADD", fmt.Sprintf("addr:%v", ntxo.Addr), bl.BlockTime, ntxokey)
            }
            log.Printf("TXOS done")

            log.Println("Before write batch LevelDB")
            err := ldb.Write(wo, wb)
            if err != nil {
                log.Fatalf("Err write batch: %v", err)
            }
            wb.Clear()
            log.Println("After write batch LevelDB")
            log.Println("Before write batch Redis")
            //wg.Add(1)
            //sem <-true
            //go func(batchtxos []interface{}) {
            //    defer wg.Done()
            //    defer func() { <-sem }()
            //db.C("txos").Insert(batchtxos...)
            //}(batchtxos)
            //r, err := conn.Do("EXEC")
            //if err != nil {
            //    panic(err)
            //}
            log.Printf("After write batch Redis")

            log.Println("Starting TXIS")
            txis := []*TxIn{}
            // Skip the ins if it's a CoinBase Tx (1 TxIn for newly generated coins)
            if !(len(tx.TxIns) == 1 && tx.TxIns[0].InputVout==0xffffffff)  {
                
                //conn.Send("MULTI")

                for txi_index, txi := range tx.TxIns {
                    ntxi := new(TxIn)
                    ntxi.TxHash = tx.Hash
                    ntxi.BlockHash = bl.Hash
                    ntxi.BlockTime = bl.BlockTime
                    ntxi.Index = uint32(txi_index)
                    nprevout := new(PrevOut)
                    nprevout.Vout = txi.InputVout
                    nprevout.Hash = txi.InputHash
                    ntxi.PrevOut = nprevout
                    prevtxo := new(TxOutCached)     
                    
                    log.Println("Starting TXI prevtxo")
                    prevtxocachedraw, err := ldb.Get(ro, []byte(fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))
                    if err != nil {
                        log.Fatalf("Err getting prevtxocached: %v", err)
                    }

                    if len(prevtxocachedraw) > 0 {
                        if err := json.Unmarshal(prevtxocachedraw, prevtxo); err != nil {
                            panic(err)
                        }
                    } else {
                        prevtxoredisjson, err := redis.String(conn.Do("GET", fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))
                        if err != nil {
                            panic(err)
                        }
                        prevtxoredis := new(TxOut)
                        json.Unmarshal([]byte(prevtxoredisjson), prevtxoredis)

                        // If something  goes wrong with LevelDB, no problem, we query MongoDB
                        //log.Println("Fallback to MongoDB")
                        //prevtxomongo := new(TxOut)
                        //if err := db.C("txos").Find(bson.M{"txhash":txi.InputHash, "index": txi.InputVout}).One(prevtxomongo); err != nil {
                        //    log.Printf("TXO requested as prevtxo: %v\n", txi.InputHash)
                        //    panic(err)                            
                        //}
                        prevtxo.Addr = prevtxoredis.Addr
                        prevtxo.Value = prevtxoredis.Value
                        //prevtxo.Id = prevtxomongo.Id.Hex()
                    }

                    log.Println("After TXI prevtxo")

                    go func(txi *blkparser.TxIn) {
                        ldb.Delete(wo, []byte(fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))
                    }(txi)
                    //} else {
                    //    for i := 1; i < 11; i++ {
                    //        err = db.C("txos").Find(bson.M{"txhash":txi.InputHash, "index": txi.InputVout}).One(prevtxo)
                    //        if err != nil {
                    //            if i == 10 {
                    //                panic(fmt.Sprintf("Can't find previous TXO for TXI: %+v, err:%v", txi, err))
                    //            }
                    //            log.Printf("Can't find previous TXO for TXI: %+v, err:%v\n", txi, err)
                    //            time.Sleep(time.Duration(i*5000)*time.Millisecond)
                    //            continue
                    //        }   
                    //    }
                    //}

                    nprevout.Address = prevtxo.Addr
                    nprevout.Value = prevtxo.Value
                    

                    txospent := new(TxoSpent)
                    txospent.Spent = true
                    txospent.BlockHeight = uint32(block_height)
                    txospent.InputHash = tx.Hash
                    txospent.InputIndex = uint32(txi_index)
                    total_tx_in+= uint(nprevout.Value)
                    txis = append(txis, ntxi)

                    //log.Println("Starting update prev txo")
                    ntxijson, _ := json.Marshal(ntxi)
                    ntxikey := fmt.Sprintf("txi:%v:%v", tx.Hash, txi_index)
                    conn.Do("SET", ntxikey, ntxijson)
                    //conn.Send("ZADD", fmt.Sprintf("txi:%v", tx.Hash), txi_index, ntxikey)

                    txospentjson, _ := json.Marshal(txospent)
                    conn.Do("SET", fmt.Sprintf("txo:%v:%v:spent", txi.InputHash, txi.InputVout), txospentjson)

                    conn.Do("ZADD", fmt.Sprintf("addr:%v", nprevout.Address), bl.BlockTime, ntxikey)
                }
                //r, err := conn.Do("EXEC")
                //if err != nil {
                //    panic(err)
                //}
            }
            
            log.Println("TXIS Done")

            total_bl_out+= total_tx_out

            ntx := new(Tx)
            ntx.Index = uint32(tx_index)
            ntx.Hash = tx.Hash
            ntx.Size = tx.Size
            ntx.LockTime = tx.LockTime
            ntx.Version = tx.Version
            ntx.TxInCnt = uint32(len(txis))
            ntx.TxOutCnt = uint32(len(txos))
            ntx.TotalOut = uint64(total_tx_out)
            ntx.TotalIn = uint64(total_tx_in)
            ntx.BlockHash = bl.Hash
            ntx.BlockHeight = block_height
            ntx.BlockTime = bl.BlockTime

            log.Println("Before TX insert")
            ntxjson, _ := json.Marshal(ntx)
            //conn.Send("MULTI")
            ntxjsonkey := fmt.Sprintf("tx:%v", ntx.Hash)
            conn.Do("SET", ntxjsonkey, ntxjson)
            conn.Do("ZADD", fmt.Sprintf("block:%v:txs", block.Hash), tx_index, ntxjsonkey)
            //conn.Do("EXEC")
            txs = append(txs, ntx)

            log.Println("TX insert done")
        }

        block.TotalBTC = uint64(total_bl_out)
        block.TxCnt = uint32(len(txs))

        log.Println("Starting block insert")
        //conn.Send("MULTI")
        blockjson, _ := json.Marshal(block)
        conn.Do("SET", fmt.Sprintf("block:%v", block.Hash), blockjson)
        conn.Do("SET", fmt.Sprintf("block:height:%v", block.Height), block.Hash)
        conn.Do("SET", "height:latest", int(block_height))
        //conn.Do("EXEC")
        log.Println("After block insert")
        //time.Sleep(100 * time.Millisecond)

        if !running {
            log.Printf("Done. Stopped at height: %v.", block_height)
        }

        wg.Done()
    }
    wg.Wait()
}
