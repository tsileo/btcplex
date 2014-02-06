package btcplex

import (
	"sort"
    "github.com/garyburd/redigo/redis"
    "fmt"
    "encoding/json"
)

type By func(tx1, tx2 *Tx) bool

func (by By) Sort(txs []*Tx) {
	ps := &txSorter{
		txs: txs,
		by:  by,
	}
	sort.Sort(ps)
}

type txSorter struct {
	txs []*Tx
	by  func(tx1, tx2 *Tx) bool
}

func (s *txSorter) Len() int {
	return len(s.txs)
}

func (s *txSorter) Swap(i, j int) {
	s.txs[i], s.txs[j] = s.txs[j], s.txs[i]
}

func (s *txSorter) Less(i, j int) bool {
	return s.by(s.txs[i], s.txs[j])
}

func TxBlockTime(tx1, tx2 *Tx) bool {
	return tx1.BlockTime < tx2.BlockTime
}

func TxFirstSeenAsc(tx1, tx2 *Tx) bool {
	return tx1.FirstSeenTime < tx2.FirstSeenTime
}

func TxFirstSeenDesc(tx1, tx2 *Tx) bool {
	return tx1.FirstSeenTime > tx2.FirstSeenTime
}

func GetUnconfirmedTxs(pool *redis.Pool) (utxs []*Tx, err error) {
	c := pool.Get()
	defer c.Close()
	utxs = []*Tx{}
	utxsid, _ := redis.Strings(c.Do("SMEMBERS", "btcplex:rawmempool"))
	for _, utxid := range utxsid {
		txraw, _ := redis.String(c.Do("GET", fmt.Sprintf("btcplex:utx:%v", utxid)))
		utx := new(Tx)
		if txraw != "" {
			err = json.Unmarshal([]byte(txraw), utx)
			if err != nil {
				return
			}
			utxs = append(utxs, utx)
		}
	}
	By(TxFirstSeenDesc).Sort(utxs)
	return
}
