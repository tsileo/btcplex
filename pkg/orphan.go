package btcplex

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func (tx *Tx) Revert(spool *redis.Pool) (err error) {
	ssdb := spool.Get()
	defer ssdb.Close()
	txupdated, err := GetTx(spool, tx.Hash)
	if err != nil {
		return
	}
	// Check if the tx is included in another block
	if tx.BlockHash == txupdated.BlockHash {
		// This tx hasn't been included in another block yet,
		// So we remove it, it will be reinserted if the tx is included
		// in a future block
		for _, addr := range tx.Addresses() {
			ssdb.Do("ZREM", fmt.Sprintf("addr:%v", addr), tx.Hash)
			ssdb.Do("ZREM", fmt.Sprintf("addr:%v:received", addr), tx.Hash)
			ssdb.Do("ZREM", fmt.Sprintf("addr:%v:sent", addr), tx.Hash)
		}
	}

	for _, txi := range tx.TxIns {
		ssdb.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txi.PrevOut.Address), "ts", -int(txi.PrevOut.Value))
	}
	for _, txo := range tx.TxOuts {
		ssdb.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txo.Addr), "tr", -int(txo.Value))
	}
	return
}
