package btcplex

import (
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
