package btcplex

import (
	"log"
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

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
            go func(conf *Config, hash string, rpool *redis.Pool, spool *redis.Pool) {
            	log.Printf("Processing new block: %v\n", hash)
                c := rpool.Get()
                defer c.Close()
                newblock, err := SaveBlockFromRPC(conf, spool, hash)
                if err != nil {
                    log.Printf("Error processing new block: %v\n", err)
                } else {
                    // Once the block is processed, we can publish it as btcplex own blocknotify
	                c.Do("PUBLISH", "btcplex:blocknotify2", hash)
	                newblockjson, _ := json.Marshal(newblock)
	                c.Do("PUBLISH", "btcplex:newblock", string(newblockjson))
                }
            }(conf, hash, rpool, spool)
        }
    }
}
