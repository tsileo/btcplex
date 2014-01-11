package btcplex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/pmylund/go-cache"
	"strconv"
)

type Block struct {
	Raw         []byte  `json:"-"`
	Hash        string  `json:"hash"`
	Height      uint    `json:"height"`
	Txs         []*Tx   `json:"tx,omitempty"`
	Version     uint32  `json:"ver"`
	MerkleRoot  string  `json:"mrkl_root"`
	BlockTime   uint32  `json:"time"`
	Bits        uint32  `json:"bits"`
	Nonce       uint32  `json:"nonce"`
	Size        uint32  `json:"size"`
	TxCnt       uint32  `json:"n_tx"`
	TotalBTC    uint64  `json:"total_out"`
	BlockReward float64 `json:"-"`
	Parent      string  `json:"prev_block"`
	Next        string  `json:"next_block"`
}

// Return block reward at the given height
func GetBlockReward(height uint32) uint64 {
	return 50e8 >> (height / 210000)
}

// Fetch a single block by hash
func GetBlock(db *levigo.DB, blockHash string, fetchTx bool) (block *Block, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	blockDataRaw, blockErr := db.Get(ro, []byte(fmt.Sprintf("bl!%s", blockHash)))
	if blockErr != nil {
		err = errors.New("Block not found")
		return
	}

	//blockStatus, _ := db.Get(ro, []byte(fmt.Sprintf("bl!%v!main", blockHash)))
	//blockPrv, _ := db.Get(ro, []byte(fmt.Sprintf("bl!%v!prv", blockHash)))

	if len(blockDataRaw) == 0 {
		err = errors.New("Block not found")
		return
	}

	block = new(Block)
	err = json.Unmarshal(blockDataRaw, block)
	if err != nil {
		return
	}

	if fetchTx {

		blockNxts, _ := GetRange(db, []byte(fmt.Sprintf("bl!%v!nxt!", block.Hash)), []byte(fmt.Sprintf("bl!%v!nxt!\xff", block.Hash)))
		for _, nxt := range blockNxts {
			blStatus, _ := db.Get(ro, []byte(fmt.Sprintf("bl!%v!main", nxt.Value)))
			blStatus2, _ := strconv.Atoi(string(blStatus[:]))
			if blStatus2 == 1 {
				block.Next = nxt.Value
			}
		}
		txs_kv, _ := GetRange(db, []byte(fmt.Sprintf("bl!%v!tx!", block.Hash)), []byte(fmt.Sprintf("bl!%v!tx!\xff", block.Hash)))

		for _, tx_kv := range txs_kv {
			tx, _ := GetTxFromKv(tx_kv)
			block.Txs = append(block.Txs, tx)
		}
	}

	return
}

// Return the Block at the given height
func GetBlockByHeight(db *levigo.DB, blockHeight uint, fetchTx bool) (block *Block, err error) {
	hash, err := GetBlockHashByHeight(db, blockHeight)
	if err != nil {
		return
	}
	block, _ = GetBlock(db, hash, fetchTx)
	return
}

// Return block hash in the main chain given the height
func GetBlockHashByHeight(db *levigo.DB, blockHeight uint) (hash string, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	blocks, _ := GetRange(db, []byte(fmt.Sprintf("bl!height!%v!", blockHeight)), []byte(fmt.Sprintf("bl!height!%v!\xff", blockHeight)))
	for _, bl := range blocks {
		blStatus, _ := db.Get(ro, []byte(fmt.Sprintf("bl!%v!main", bl.Value)))
		blStatus2, _ := strconv.Atoi(string(blStatus[:]))
		if blStatus2 == 1 {
			hash = bl.Value
			return
		}
	}
	err = errors.New("Block not found")
	return
}

func GetBlockCached(c *cache.Cache, db *levigo.DB, blockHash string, fetchTx bool) (block *Block, err error) {
	cachekey := fmt.Sprintf("block%v", blockHash)
	blockcached, found := c.Get(cachekey)
	if found {
		return blockcached.(*Block), nil
	} else {
		block, err := GetBlock(db, blockHash, fetchTx)
		c.Set(cachekey, block, 0)
		return block, err
	}
}
