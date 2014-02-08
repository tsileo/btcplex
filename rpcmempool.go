package btcplex

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"time"
    "github.com/garyburd/redigo/redis"
)

// Get unconfirmed transactions from memory pool, along with
// first seem time/block height, requires a recent bitcoind version
func ProcessUnconfirmedTxsRPC(conf *Config, pool *redis.Pool, running *bool) {    
	c := pool.Get()
	defer c.Close()
	c.Do("DEL", "btcplex:rawmempool")
	lastkey := ""
	ckey := ""
	for {
		if !*running {
        	log.Println("Stopping ProcessUnconfirmedTxsRPC")
            break
        }
        ckey = fmt.Sprintf("btcplex:rawmempool:%v", time.Now().UTC().Unix())
        log.Printf("lastkey:%+v, ckey:%+v\n", lastkey, ckey)
       	unconfirmedtxsverbose, _ := GetRawMemPoolVerboseRPC(conf)
		unconfirmedtxs, _ := GetRawMemPoolRPC(conf)
		for _, txid := range unconfirmedtxs {
			fmt.Printf("UTX:%v\n", txid)
			txkey := fmt.Sprintf("btcplex:utx:%v", txid)
			txexists, _ := redis.Bool(c.Do("EXISTS", txkey))
			txmeta, txmetafound := unconfirmedtxsverbose[txid].(map[string]interface{})
			if txmetafound {

				fseentime, _ := txmeta["time"].(json.Number).Int64()
				if !txexists {
					tx, _ := GetTxRPC(conf, txid, &Block{})
					tx.FirstSeenTime = uint32(fseentime)
					fseenheight, _ := txmeta["height"].(json.Number).Int64()
					tx.FirstSeenHeight = uint(fseenheight)
					txjson, _ := json.Marshal(tx)
					c.Do("SET", txkey, string(txjson))
					c.Do("PUBLISH", "btcplex:utxs", string(txjson))
				}
				c.Do("ZADD", "btcplex:rawmempool", fseentime, txkey)
				c.Do("SADD", ckey, txkey)	
			} else {
				log.Printf("Error utx:%v", txid)
			}
		}
        if lastkey == "" {
        	lastkey = ckey
        } else {
        	dkeys, _ := redis.Strings(c.Do("SDIFF", lastkey, ckey))
        	log.Printf("Deleting %v utxs\n", len(dkeys))
        	c.Do("DEL", redis.Args{}.Add(lastkey).AddFlat(dkeys)...)
        	lastkey = ckey
        	c.Do("ZREM", redis.Args{}.Add("btcplex:rawmempool").AddFlat(dkeys)...)
        }
		time.Sleep(1 * time.Second)
	}
}
