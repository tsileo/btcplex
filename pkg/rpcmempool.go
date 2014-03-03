package btcplex

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "io/ioutil"
	"log"
	"sync"
	"time"
)

// Get unconfirmed transactions from memory pool, along with
// first seem time/block height, requires a recent bitcoind version
func ProcessUnconfirmedTxs(conf *Config, pool *redis.Pool, running *bool) {
	var wg sync.WaitGroup
	var lastts, cts int64
	var lastkey, ckey string

	log.Println("ProcessUnconfirmedTxs startup")

	c := pool.Get()
	defer c.Close()

	// Cleanup old keys since it has stopped
	oldkeys, _ := redis.Strings(c.Do("ZRANGE", "btcplex:rawmempool", 0, -1))
	c.Do("DEL", redis.Args{}.AddFlat(oldkeys)...)
	c.Do("DEL", "btcplex:rawmempool")

	// We fetch 25 tx max in the pool
	sem := make(chan bool, 25)

	for {
		if !*running {
			log.Println("Stopping ProcessUnconfirmedTxs")
			break
		}

		cts = time.Now().UTC().Unix()
		ckey = fmt.Sprintf("btcplex:rawmempool:%v", cts)

		//log.Printf("lastkey:%+v, ckey:%+v\n", lastkey, ckey)

		// Call bitcoind RPC
		unconfirmedtxsverbose, _ := GetRawMemPoolVerboseRPC(conf)
		unconfirmedtxs, _ := GetRawMemPoolRPC(conf)

		for _, txid := range unconfirmedtxs {
			wg.Add(1)
			sem <- true
			go func(pool *redis.Pool, txid string, unconfirmedtxsverbose *map[string]interface{}) {
				c := pool.Get()
				defer c.Close()
				defer wg.Done()
				defer func() { <-sem }()
				txkey := fmt.Sprintf("btcplex:utx:%v", txid)
				txexists, _ := redis.Bool(c.Do("EXISTS", txkey))
				uverbose := *unconfirmedtxsverbose
				txmeta, txmetafound := uverbose[txid].(map[string]interface{})
				if txmetafound {
					fseentime, _ := txmeta["time"].(json.Number).Int64()
					if !txexists {
						tx, _ := GetTxRPC(conf, txid, &Block{})
						tx.FirstSeenTime = uint32(fseentime)
						fseenheight, _ := txmeta["height"].(json.Number).Int64()
						tx.FirstSeenHeight = uint(fseenheight)
						txjson, _ := json.Marshal(tx)
						c.Do("SET", txkey, string(txjson))
					}
					c.Do("ZADD", "btcplex:rawmempool", fseentime, txkey)
					// Put the TX in a snapshot do detect deleted tx
					c.Do("SADD", ckey, txkey)
				}
			}(pool, txid, &unconfirmedtxsverbose)
		}
		wg.Wait()
		if lastkey != "" {
			// We remove tx that are no longer in the pool using the last snapshot
			dkeys, _ := redis.Strings(c.Do("SDIFF", lastkey, ckey))
			//log.Printf("Deleting %v utxs\n", len(dkeys))
			c.Do("DEL", redis.Args{}.Add(lastkey).AddFlat(dkeys)...)
			c.Do("ZREM", redis.Args{}.Add("btcplex:rawmempool").AddFlat(dkeys)...)
			// Since getrawmempool return transaction sorted by name, we replay them sorted by time asc
			newkeys, _ := redis.Strings(c.Do("ZRANGEBYSCORE", "btcplex:rawmempool", fmt.Sprintf("(%v", lastts), cts))
			for _, newkey := range newkeys {
				txjson, _ := redis.String(c.Do("GET", newkey))
				// Notify SSE unconfirmed transactions
				c.Do("PUBLISH", "btcplex:utxs", txjson)
				ctx := new(Tx)
				json.Unmarshal([]byte(txjson), ctx)
				// Notify transaction to every channel address
				multiPublishScript.Do(c, redis.Args{}.Add(txjson).AddFlat(ctx.AddressesChannels())...)
				c.Do("SETEX", fmt.Sprintf("btcplex:utx:%v:published", ctx.Hash), 3600*20, cts)
				//c.Do("SADD", "btcplex:utxs:published", ctx.Hash)
			}
		} else {
			log.Println("ProcessUnconfirmedTxs first round done")
		}
		lastkey = ckey
		lastts = cts
		time.Sleep(500 * time.Millisecond)
	}
}

// Fetch unconfirmed tx from Redis
func GetUnconfirmedTx(pool *redis.Pool, hash string) (tx *Tx, err error) {
	c := pool.Get()
	defer c.Close()
	txkey := fmt.Sprintf("btcplex:utx:%v", hash)
	tx = new(Tx)
	txjson, _ := redis.String(c.Do("GET", txkey))
	err = json.Unmarshal([]byte(txjson), tx)
	return
}
