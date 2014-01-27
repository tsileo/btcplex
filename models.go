package btcplex

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
)

type Block struct {
	Id         bson.ObjectId `json:"-" bson:"_id"`
	Hash       string        `json:"hash"`
	Height     uint          `json:"height"`
	Txs        []*Tx         `json:"tx,omitempty" bson:"-"`
	Version    uint32        `json:"ver"`
	MerkleRoot string        `json:"mrkl_root"`
	BlockTime  uint32        `json:"time"`
	Bits       uint32        `json:"bits"`
	Nonce      uint32        `json:"nonce"`
	Size       uint32        `json:"size"`
	TxCnt      uint32        `json:"n_tx"`
	TotalBTC   uint64        `json:"total_out"`
	//    BlockReward float64 `json:"-"`
	Parent string `json:"prev_block"`
	Next   string `json:"next_block"`
}

type Tx struct {
	Id              bson.ObjectId  `json:"-" bson:"_id"`
	Hash            string         `json:"hash"`
	Index           uint32         `json:"-"`
	Size            uint32         `json:"size"`
	LockTime        uint32         `json:"lock_time"`
	Version         uint32         `json:"ver"`
	TxInCnt         uint32         `json:"vin_sz"`
	TxOutCnt        uint32         `json:"vout_sz"`
	TxIns           []*TxIn        `json:"in" bson:"-"`
	TxOuts          []*TxOut       `json:"out" bson:"-"`
	TotalOut        uint64         `json:"vout_total"`
	TotalIn         uint64         `json:"vin_total"`
	BlockHash       string         `json:"block_hash"`
	BlockHeight     uint           `json:"block_height"`
	BlockTime       uint32         `json:"block_time"`
	FirstSeenTime   uint32         `json:"-"`
	FirstSeenHeight uint           `json:"-"`
	TxAddressInfo   *TxAddressInfo `json:"-"`
}

type TxAddressInfo struct {
	InTxIn  bool
	InTxOut bool
	Value   int64
}

type TxOut struct {
	Id        bson.ObjectId `json:"-" bson:"_id"`
	TxHash    string        `json:"-"`
	BlockHash string        `json:"-"`
	BlockTime uint32        `json:"-"`
	Addr      string        `json:"hash"`
	Value     uint64        `json:"value"`
	Index     uint32        `json:"n"`
	Spent     *TxoSpent     `json:"spent,omitempty"`
}

type PrevOut struct {
	Hash    string `json:"hash"`
	Vout    uint32 `json:"n"`
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

type TxIn struct {
	Id        bson.ObjectId `json:"-" bson:"_id"`
	TxHash    string        `json:"-"`
	BlockHash string        `json:"-"`
	BlockTime uint32        `json:"-"`
	PrevOut   *PrevOut      `json:"prev_out"`
	Index     uint32        `json:"n"`
}

type TxoSpent struct {
	Spent       bool   `json:"spent"`
	BlockHeight uint32 `json:"block_height,omitempty" gorethink:",omitempty"`
	InputHash   string `json:"tx_hash,omitempty"  gorethink:",omitempty"`
	InputIndex  uint32 `json:"in_index,omitempty" gorethink:",omitempty"`
}

// Return block reward at the given height
func GetBlockReward(height uint) uint {
	return 50e8 >> (height / 210000)
}

// Fetch a block by height
func GetBlockByHeight(db *mgo.Database, height uint) (block *Block, err error) {
	block = new(Block)
	err = db.C("blocks").Find(bson.M{"height": height}).One(block)
	if err != nil {
		return
	}
	return
}

// Fetch a block by hash
func GetBlockByHash(db *mgo.Database, hash string) (block *Block, err error) {
	block = new(Block)
	err = db.C("blocks").Find(bson.M{"hash": hash}).One(block)
	if err != nil {
		return
	}
	return
}

// Fetch a transaction by hash
func GetTx(db *mgo.Database, hash string) (tx *Tx, err error) {
	tx = new(Tx)
	err = db.C("txs").Find(bson.M{"hash": hash}).One(tx)
	if err != nil {
		return
	}
	return
}

// Fetch Txos and Txins
func (tx *Tx) Build(db *mgo.Database) (err error) {
	tx.TxIns = []*TxIn{}
	err = db.C("txis").Find(bson.M{"txhash": tx.Hash}).Sort("index").All(&tx.TxIns)
	if err != nil {
		return
	}
	tx.TxOuts = []*TxOut{}
	err = db.C("txos").Find(bson.M{"txhash": tx.Hash}).Sort("index").All(&tx.TxOuts)
	if err != nil {
		return
	}
	return
}

// Fetch all block transactions
func (block *Block) FetchTxs(db *mgo.Database) (err error) {
	block.Txs = []*Tx{}
	err = db.C("txs").Find(bson.M{"blockhash": block.Hash}).Sort("index").All(&block.Txs)
	if err != nil {
		return
	}
	for _, tx := range block.Txs {
		tx.Build(db)
	}
	return
}

// Return last X blocks from stop to start (both included)
func GetLastXBlocks(db *mgo.Database, start uint, stop uint) (blocks []*Block, err error) {
	blocks = []*Block{}
	err = db.C("blocks").Find(bson.M{"height": bson.M{"$lte": start, "$gte": stop}}).Sort("-height").All(&blocks)
	if err != nil {
		return
	}
	return
}
