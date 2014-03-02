package btcplex

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/iter"
	"github.com/garyburd/redigo/redis"
)

type Block struct {
	Hash       string `json:"hash"`
	Height     uint   `json:"height"`
	Txs        []*Tx  `json:"tx,omitempty"`
	Version    uint32 `json:"ver"`
	MerkleRoot string `json:"mrkl_root"`
	BlockTime  uint32 `json:"time"`
	Bits       uint32 `json:"bits"`
	Nonce      uint32 `json:"nonce"`
	Size       uint32 `json:"size"`
	TxCnt      uint32 `json:"n_tx"`
	TotalBTC   uint64 `json:"total_out"`
	//    BlockReward float64 `json:"-"`
	Parent string                       `json:"prev_block"`
	Next   string                       `json:"next_block"`
	Links  map[string]map[string]string `json:"_links,omitempty"`
	Meta            *BlockMeta                   `json:"-"`
	Main bool `json:"-"`
}

type Tx struct {
	Hash            string                       `json:"hash"`
	Index           uint32                       `json:"-"`
	Size            uint32                       `json:"size"`
	LockTime        uint32                       `json:"lock_time"`
	Version         uint32                       `json:"ver"`
	TxInCnt         uint32                       `json:"vin_sz"`
	TxOutCnt        uint32                       `json:"vout_sz"`
	TxIns           []*TxIn                      `json:"in"`
	TxOuts          []*TxOut                     `json:"out"`
	TotalOut        uint64                       `json:"vout_total"`
	TotalIn         uint64                       `json:"vin_total"`
	BlockHash       string                       `json:"block_hash"`
	BlockHeight     uint                         `json:"block_height"`
	BlockTime       uint32                       `json:"block_time"`
	FirstSeenTime   uint32                       `json:"first_seen_time"`
	FirstSeenHeight uint                         `json:"first_seen_height"`
	TxAddressInfo   *TxAddressInfo               `json:"-"`
	Links           map[string]map[string]string `json:"_links,omitempty"`
}

type TxAddressInfo struct {
	InTxIn  bool
	InTxOut bool
	Value   int64
}

type TxOut struct {
	TxHash    string    `json:"-"`
	BlockHash string    `json:"-"`
	BlockTime uint32    `json:"-"`
	Addr      string    `json:"hash"`
	Value     uint64    `json:"value"`
	Index     uint32    `json:"n"`
	Spent     *TxoSpent `json:"spent,omitempty"`
}

type PrevOut struct {
	Hash    string `json:"hash"`
	Vout    uint32 `json:"n"`
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

type TxIn struct {
	TxHash    string   `json:"-"`
	BlockHash string   `json:"-"`
	BlockTime uint32   `json:"-"`
	PrevOut   *PrevOut `json:"prev_out"`
	Index     uint32   `json:"n"`
}

type TxoSpent struct {
	Spent       bool   `json:"spent"`
	BlockHeight uint32 `json:"block_height,omitempty"`
	InputHash   string `json:"tx_hash,omitempty"`
	InputIndex  uint32 `json:"in_index,omitempty"`
}

type BlockMeta struct {
	Main	bool	`redis:"main"`
	Next    string  `redis:"next"`
	Parent  string  `redis:"parent"`
	Height  int     `redis:"height"`  
}

// Return block reward at the given height
func GetBlockReward(height uint) uint {
	return 50e8 >> (height / 210000)
}

// Return block hash for the given height
func GetBlockHash(rpool *redis.Pool, height uint) (hash string, err error) {
	c := rpool.Get()
	defer c.Close()
	hash, err = redis.String(c.Do("GET", fmt.Sprintf("block:height:%v", height)))
	return
}

// Get a block by its hash
func GetBlockByHash(rpool *redis.Pool, hash string) (block *Block, err error) {
	c := rpool.Get()
	defer c.Close()
	blockjson, err := redis.String(c.Do("GET", fmt.Sprintf("block:%v", hash)))
	if err != nil {
		return
	}
	block = new(Block)
	err = json.Unmarshal([]byte(blockjson), block)
	return
}

// Get a block by its hash along with its full transactions
func GetBlockCachedByHash(rpool *redis.Pool, hash string) (block *Block, err error) {
	c := rpool.Get()
	defer c.Close()
	blockjson, err := redis.String(c.Do("GET", fmt.Sprintf("block:%v:cached", hash)))
	if err != nil {
		return
	}
	block = new(Block)
	err = json.Unmarshal([]byte(blockjson), block)
	return
}

// TODO UpdateTxoSpent

func (block *Block) FetchTxs(rpool *redis.Pool) (err error) {
	c := rpool.Get()
	defer c.Close()
	txskeys, _ := redis.Strings(c.Do("ZRANGE", fmt.Sprintf("block:%v:txs", block.Hash), 0, 1000))
	txskeysi := []interface{}{}
	for _, txkey := range txskeys {
		txskeysi = append(txskeysi, txkey)
	}
	txsjson, _ := redis.Strings(c.Do("MGET", txskeysi...))
	block.Txs = []*Tx{}
	for _, txjson := range txsjson {
		ctx := new(Tx)
		err = json.Unmarshal([]byte(txjson), ctx)
		if err != nil {
			return
		}
		ctx.Build(rpool)
		block.Txs = append(block.Txs, ctx)
	}
	return
}

func (block *Block) FetchMeta(rpool *redis.Pool) (err error) {
	c := rpool.Get()
	defer c.Close()
	meta, err := NewBlockMeta(rpool, block.Hash)
	if err != nil {
		return
	}
	block.Next = meta.Next
	block.Main = meta.Main
	return
}

func NewBlockMeta(rpool *redis.Pool, block_hash string) (blockmeta *BlockMeta, err error) {
	c := rpool.Get()
	defer c.Close()
	blockmeta = new(BlockMeta)
	v, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("block:%v:h", block_hash)))
	if err != nil {
		return
	}
	if err = redis.ScanStruct(v, blockmeta); err != nil {
		return
	}
	return
}

// Fetch a transaction by hash
func GetTx(rpool *redis.Pool, hash string) (tx *Tx, err error) {
	c := rpool.Get()
	defer c.Close()
	tx = new(Tx)
	txjson, _ := redis.String(c.Do("GET", fmt.Sprintf("tx:%v", hash)))
	err = json.Unmarshal([]byte(txjson), tx)
	tx.Build(rpool)
	return
}

// Fetch Txos and Txins
func (tx *Tx) Build(rpool *redis.Pool) (err error) {
	c := rpool.Get()
	defer c.Close()
	tx.TxIns = []*TxIn{}
	tx.TxOuts = []*TxOut{}
	txinskeys := []interface{}{}
	for i := range iter.N(int(tx.TxInCnt)) {
		txinskeys = append(txinskeys, fmt.Sprintf("txi:%v:%v", tx.Hash, i))
	}
	txinsjson, _ := redis.Strings(c.Do("MGET", txinskeys...))
	for _, txinjson := range txinsjson {
		ctxi := new(TxIn)
		err = json.Unmarshal([]byte(txinjson), ctxi)
		tx.TxIns = append(tx.TxIns, ctxi)
	}
	txoutskeys := []interface{}{}
	txoutsspentkeys := []interface{}{}
	for i := range iter.N(int(tx.TxOutCnt)) {
		txoutskeys = append(txoutskeys, fmt.Sprintf("txo:%v:%v", tx.Hash, i))
		txoutsspentkeys = append(txoutsspentkeys, fmt.Sprintf("txo:%v:%v:spent", tx.Hash, i))
	}
	txoutsjson, _ := redis.Strings(c.Do("MGET", txoutskeys...))
	txoutsspentjson, _ := redis.Strings(c.Do("MGET", txoutsspentkeys...))
	for txoindex, txoutjson := range txoutsjson {
		ctxo := new(TxOut)
		err = json.Unmarshal([]byte(txoutjson), ctxo)
		if txoutsspentjson[txoindex] != "" {
			cspent := new(TxoSpent)
			err = json.Unmarshal([]byte(txoutsspentjson[txoindex]), cspent)
			ctxo.Spent = cspent
		}
		tx.TxOuts = append(tx.TxOuts, ctxo)
	}
	return
}

// Return last X blocks from stop to start (both included)
func GetLastXBlocks(rpool *redis.Pool, start uint, stop uint) (blocks []*Block, err error) {
	c := rpool.Get()
	defer c.Close()
	blocks = []*Block{}
	cur := int(start)
	blockskeys := []interface{}{}
	for _ = range iter.N(int(start - stop)) {
		// TODO(tsileo) MGET here to retrieve hashes
		chash, cerr := GetBlockHash(rpool, uint(cur))
		if cerr != nil {
			err = cerr
			return
		}
		blockskeys = append(blockskeys, fmt.Sprintf("block:%v", chash))
		cur -= 1
	}
	blocksjson, _ := redis.Strings(c.Do("MGET", blockskeys...))
	for _, blockjson := range blocksjson {
		cblock := new(Block)
		err = json.Unmarshal([]byte(blockjson), cblock)
		blocks = append(blocks, cblock)
	}
	return
}
