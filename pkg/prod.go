package btcplex

import (
	"log"
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

func CatchUpLatestBlock(conf *Config, rpool *redis.Pool, spool *redis.Pool) (done bool) {
    blockcount := GetBlockCountRPC(conf)
    sc := spool.Get()
    defer sc.Close()
    latestheight, _ := redis.Int(sc.Do("GET", "height:latest"))
    if uint(latestheight) == blockcount {
        return true
    }
    hash := GetBlockHashRPC(conf, uint(latestheight) + 1)
    log.Printf("Catch up block: %v\n", hash)
    SaveBlockFromRPC(conf, spool, hash)
    return false
}

func ProcessNewBlock(conf *Config, rpool *redis.Pool, spool *redis.Pool) {
	log.Println("ProcessNewBlock startup")
    conn := rpool.Get()
    defer conn.Close()
    psc := redis.PubSubConn{Conn: conn}
    psc.Subscribe("btcplex:blocknotify")
    for {
        switch v := psc.Receive().(type) {
        case redis.Message:
            hash := string(v.Data)
            log.Printf("Processing new block: %v\n", hash)
            c := rpool.Get()
            newblock, err := SaveBlockFromRPC(conf, spool, hash)
            if err != nil {
                log.Printf("Error processing new block: %v\n", err)
            } else {
                // Once the block is processed, we can publish it as btcplex own blocknotify
	            c.Do("PUBLISH", "btcplex:blocknotify2", hash)
	            newblockjson, _ := json.Marshal(newblock)
	            c.Do("PUBLISH", "btcplex:newblock", string(newblockjson))
            }
            c.Close()
        }
    }
}
