package btcplex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/pmylund/go-cache"
	"strconv"
	"strings"
)

type Tx struct {
	Hash          string         `json:"hash"`
	Size          uint32         `json:"size"`
	LockTime      uint32         `json:"lock_time"`
	Version       uint32         `json:"ver"`
	TxInCnt       uint32         `json:"vin_sz"`
	TxOutCnt      uint32         `json:"vout_sz"`
	TxIns         []*TxIn        `json:"in"`
	TxOuts        []*TxOut       `json:"out"`
	TotalOut      uint64         `json:"vout_total"`
	TotalIn       uint64         `json:"vin_total"`
	BlockHash     string         `json:"block_hash"`
	BlockHeight   uint           `json:"block_height"`
	BlockTime     uint32         `json:"block_time"`
	TxAddressInfo *TxAddressInfo `json:"-"`
}

type TxAddressInfo struct {
	InTxIn  bool
	InTxOut bool
	Value   int64
}

type PrevOut struct {
	Hash    string `json:"hash"`
	Vout    uint32 `json:"n"`
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

type TxOut struct {
	Addr     string    `json:"hash"`
	Value    uint64    `json:"value"`
	Pkscript string    `json:"-"`
	Spent    *TxoSpent `json:"spent,omitempty"`
}

type TxIn struct {
	PrevOut   *PrevOut `json:"prev_out"`
	ScriptSig string   `json:"-"`
}

type TxoSpent struct {
	Spent       bool   `json:"spent"`
	BlockHeight uint32 `json:"block_height,omitempty"`
	InputHash   string `json:"tx_hash,omitempty"`
	InputIndex  uint32 `json:"in_index,omitempty"`
}

// Fetch a Tx by hash
func GetTx(db *levigo.DB, txHash string) (tx *Tx, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	txKey, blockErr := db.Get(ro, []byte(fmt.Sprintf("tx!%s", txHash)))
	if blockErr != nil {
		err = errors.New("Tx not found")
		return
	}

	txData, blockErr := db.Get(ro, txKey)
	if blockErr != nil {
		err = errors.New("Tx not found")
		return
	}

	tx = new(Tx)
	err = json.Unmarshal(txData, tx)
	if err != nil {
		return
	}

	for txo_index, txo := range tx.TxOuts {
		txo.Spent, _ = GetTxoSpent(db, txo.Addr, tx.Hash, txo_index)
	}

	return
}

func GetTxoSpent(db *levigo.DB, address string, txhhash string, txoindex int) (txospent *TxoSpent, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	txSpent, err := db.Get(ro, []byte(fmt.Sprintf("%v-txo-spent!%v!%v", address, txhhash, txoindex)))
	if err != nil {
		return
	}
	txospent = new(TxoSpent)
	if string(txSpent[:]) == "0" || string(txSpent[:]) == "" {
		txospent.Spent = false
		return
	}
	txospent.Spent = true

	txSpentData := strings.Split(string(txSpent[:]), ":")
	blockHeight, _ := strconv.ParseUint(txSpentData[0], 10, 0)
	txiIndex, _ := strconv.ParseUint(txSpentData[2], 10, 0)

	txospent.BlockHeight = uint32(blockHeight)
	txospent.InputHash = txSpentData[1]
	txospent.InputIndex = uint32(txiIndex)

	return
}

// Fetch a Tx by hash
func GetTxFromKv(kv *KeyValue) (tx *Tx, err error) {
	tx = new(Tx)
	err = json.Unmarshal([]byte(kv.Value), tx)
	return
}

func GetTxCached(c *cache.Cache, db *levigo.DB, txHash string) (tx *Tx, err error) {
	cachekey := fmt.Sprintf("tx%v", txHash)
	txcached, found := c.Get(cachekey)
	if found {
		return txcached.(*Tx), nil
	} else {
		tx, err := GetTx(db, txHash)
		c.Set(cachekey, tx, 0)
		return tx, err
	}
}
