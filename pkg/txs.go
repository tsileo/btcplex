package btcplex

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sort"
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

func TxIndex(tx1, tx2 *Tx) bool {
	return tx1.Index < tx2.Index
}

// Return all unconfirmed transactions from Redis
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

// Return a set containing every addresses listed in txis/txos
func (tx *Tx) Addresses() (addresses []string) {
	addrset := make(map[string]struct{})
	for _, txi := range tx.TxIns {
		addrset[txi.PrevOut.Address] = struct{}{}
	}
	for _, txo := range tx.TxOuts {
		addrset[txo.Addr] = struct{}{}
	}

	addresses = []string{}
	for addr, _ := range addrset {
		addresses = append(addresses, addr)
	}
	return
}

// Return a set containing every addresses listed in txis/txos
func (tx *Tx) AddressesChannels() (addresses []string) {
	addrset := make(map[string]struct{})
	for _, txi := range tx.TxIns {
		addrset[txi.PrevOut.Address] = struct{}{}
	}
	for _, txo := range tx.TxOuts {
		addrset[txo.Addr] = struct{}{}
	}

	addresses = []string{}
	for addr, _ := range addrset {
		addresses = append(addresses, fmt.Sprintf("addr:%v:txs", addr))
	}
	return
}
