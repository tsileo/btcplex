package btcplex

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func (tx *Tx) Revert(spool *redis.Pool) (err error) {
	ssdb := spool.Get()
	defer ssdb.Close()

	for _, txi := range tx.TxIns {
		ssdb.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txi.PrevOut.Address), "ts", -int(txi.PrevOut.Value))
		ssdb.Do("ZREM", fmt.Sprintf("addr:%v", txi.PrevOut.Address), tx.Hash)
		ssdb.Do("ZREM", fmt.Sprintf("addr:%v:sent", txi.PrevOut.Address), tx.Hash)

	}
	for _, txo := range tx.TxOuts {
		ssdb.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txo.Addr), "tr", -int(txo.Value))
		ssdb.Do("ZREM", fmt.Sprintf("addr:%v", txo.Addr), tx.Hash)
		ssdb.Do("ZREM", fmt.Sprintf("addr:%v:received", txo.Addr), tx.Hash)
	}

	return
}
