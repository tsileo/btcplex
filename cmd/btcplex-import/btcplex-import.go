package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jmhodges/levigo"
	blkparser "github.com/tsileo/blkparser"

	"btcplex"
)

type Block struct {
	Hash       string `json:"hash"`
	Height     uint   `json:"height"`
	Txs        []*Tx  `json:"tx,omitempty" bson:"-"`
	Version    uint32 `json:"ver"`
	MerkleRoot string `json:"mrkl_root"`
	BlockTime  uint32 `json:"time"`
	Bits       uint32 `json:"bits"`
	Nonce      uint32 `json:"nonce"`
	Size       uint32 `json:"size"`
	TxCnt      uint32 `json:"n_tx"`
	TotalBTC   uint64 `json:"total_out"`
	//    BlockReward float64 `json:"-"`
	Parent string `json:"prev_block"`
	//    Next        string  `json:"next_block"`
}

type Tx struct {
	Hash        string   `json:"hash"`
	Index       uint32   `json:"-"`
	Size        uint32   `json:"size"`
	LockTime    uint32   `json:"lock_time"`
	Version     uint32   `json:"ver"`
	TxInCnt     uint32   `json:"vin_sz"`
	TxOutCnt    uint32   `json:"vout_sz"`
	TxIns       []*TxIn  `json:"in" bson:"-"`
	TxOuts      []*TxOut `json:"out" bson:"-"`
	TotalOut    uint64   `json:"vout_total"`
	TotalIn     uint64   `json:"vin_total"`
	BlockHash   string   `json:"block_hash"`
	BlockHeight uint     `json:"block_height"`
	BlockTime   uint32   `json:"block_time"`
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

type TxOutCached struct {
	Id    string `json:"id"`
	Addr  string `json:"hash"`
	Value uint64 `json:"value"`
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
	BlockHeight uint32 `json:"block_height"`
	InputHash   string `json:"tx_hash"`
	InputIndex  uint32 `json:"in_index"`
}

var wg, txwg sync.WaitGroup
var tximut, txomut sync.Mutex

var running bool

func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}

func main() {
	fmt.Printf("GOMAXPROCS is %d\n", getGOMAXPROCS())
	confFile := "config.json"
		conf, err := btcplex.LoadConfig(confFile)
	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}
	pool, err := btcplex.GetSSDB(conf)
	if err != nil {
		log.Fatalf("Can't connect to SSDB: %v", err)
	}

	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)
	ldb, err := levigo.Open(conf.LevelDbPath, opts) //alpha
	defer ldb.Close()

	if err != nil {
		log.Fatalf("failed to load db: %s\n", err)
	}

	wo := levigo.NewWriteOptions()
	//wo.SetSync(true)
	defer wo.Close()

	ro := levigo.NewReadOptions()
	defer ro.Close()

	wb := levigo.NewWriteBatch()
	defer wb.Close()

	conn := pool.Get()
	defer conn.Close()

	log.Println("Waiting 3 seconds before starting...")
	time.Sleep(3 * time.Second)

	latestheight := 0
	log.Printf("Latest height: %v\n", latestheight)

	running = true
	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt)
	go func() {
		for sig := range cs {
			running = false
			log.Printf("Captured %v, waiting for everything to finish...\n", sig)
			wg.Wait()
			defer os.Exit(1)
		}
	}()

	concurrency := 250
	sem := make(chan bool, concurrency)

	// Real network magic byte
	blockchain, blockchainerr := blkparser.NewBlockchain(conf.BitcoindBlocksPath, [4]byte{0xF9, 0xBE, 0xB4, 0xD9})
	if blockchainerr != nil {
		log.Fatalf("Error loading block file: ", blockchainerr)
	}

	block_height := uint(0)
	for {
		if !running {
			break
		}

		wg.Add(1)

		bl, er := blockchain.NextBlock()
		if er != nil {
			log.Println("Initial import done.")
			break
		}

		bl.Raw = nil

		if bl.Parent == "" {
			block_height = uint(0)
			conn.Do("HSET", fmt.Sprintf("block:%v:h", bl.Hash), "main", true)
			conn.Do("HSET", fmt.Sprintf("block:%v:h", bl.Hash), "height", 0)

		} else {
			parentheight, _ := redis.Int(conn.Do("HGET", fmt.Sprintf("block:%v:h", bl.Parent), "height"))
			block_height = uint(parentheight + 1)
			conn.Do("HSET", fmt.Sprintf("block:%v:h", bl.Hash), "height", block_height)
			prevheight := block_height - 1
			prevhashtest := bl.Parent
			prevnext := bl.Hash
			for {
				prevkey := fmt.Sprintf("height:%v", prevheight)
				prevcnt, _ := redis.Int(conn.Do("ZCARD", prevkey))
				// SSDB doesn't support negative slice yet
				prevs, _ := redis.Strings(conn.Do("ZRANGE", prevkey, 0, prevcnt-1))
				for _, cprevhash := range prevs {
					if cprevhash == prevhashtest {
						// current block parent
						prevhashtest, _ = redis.String(conn.Do("HGET", fmt.Sprintf("block:%v:h", cprevhash), "parent"))
						// Set main to 1 and the next => prevnext
						conn.Do("HMSET", fmt.Sprintf("block:%v:h", cprevhash), "main", true, "next", prevnext)
						conn.Do("SET", fmt.Sprintf("block:height:%v", prevheight), cprevhash)
						prevnext = cprevhash
					} else {
						// Set main to 0
						conn.Do("HSET", fmt.Sprintf("block:%v:h", cprevhash), "main", false)
						oblock, _ := btcplex.GetBlockCachedByHash(pool, cprevhash)
						for _, otx := range oblock.Txs {
							otx.Revert(pool)
						}
					}
				}
				if len(prevs) == 1 {
					break
				}
				prevheight--
			}
			//}

		}

		// Orphans blocks handling
		conn.Do("ZADD", fmt.Sprintf("height:%v", block_height), bl.BlockTime, bl.Hash)
		conn.Do("HSET", fmt.Sprintf("block:%v:h", bl.Hash), "parent", bl.Parent)

		if latestheight != 0 && !(latestheight+1 <= int(block_height)) {
			log.Printf("Skipping block #%v\n", block_height)
			continue
		}

		log.Printf("Current block: %v (%v)\n", block_height, bl.Hash)

		block := new(Block)
		block.Hash = bl.Hash
		block.Height = block_height
		block.Version = bl.Version
		block.MerkleRoot = bl.MerkleRoot
		block.BlockTime = bl.BlockTime
		block.Bits = bl.Bits
		block.Nonce = bl.Nonce
		block.Size = bl.Size
		block.Parent = bl.Parent

		txs := []*Tx{}

		total_bl_out := uint64(0)
		for tx_index, tx := range bl.Txs {
			//log.Printf("Tx #%v: %v\n", tx_index, tx.Hash)

			total_tx_out := uint64(0)
			total_tx_in := uint64(0)

			//conn.Send("MULTI")
			txos := []*TxOut{}
			txis := []*TxIn{}

			for txo_index, txo := range tx.TxOuts {
				txwg.Add(1)
				sem <- true
				go func(bl *blkparser.Block, tx *blkparser.Tx, pool *redis.Pool, total_tx_out *uint64, txo *blkparser.TxOut, txo_index int) {
					conn := pool.Get()
					defer conn.Close()
					defer func() {
						<-sem
					}()
					defer txwg.Done()
					atomic.AddUint64(total_tx_out, uint64(txo.Value))
					//atomic.AddUint32(txos_cnt, 1)

					ntxo := new(TxOut)
					ntxo.TxHash = tx.Hash
					ntxo.BlockHash = bl.Hash
					ntxo.BlockTime = bl.BlockTime
					ntxo.Addr = txo.Addr
					ntxo.Value = txo.Value
					ntxo.Index = uint32(txo_index)
					txospent := new(TxoSpent)
					ntxo.Spent = txospent
					ntxocached := new(TxOutCached)
					ntxocached.Addr = txo.Addr
					ntxocached.Value = txo.Value

					ntxocachedjson, _ := json.Marshal(ntxocached)
					ldb.Put(wo, []byte(fmt.Sprintf("txo:%v:%v", tx.Hash, txo_index)), ntxocachedjson)

					ntxojson, _ := json.Marshal(ntxo)
					ntxokey := fmt.Sprintf("txo:%v:%v", tx.Hash, txo_index)
					conn.Do("SET", ntxokey, ntxojson)

					//conn.Send("ZADD", fmt.Sprintf("txo:%v", tx.Hash), txo_index, ntxokey)
					conn.Do("ZADD", fmt.Sprintf("addr:%v", ntxo.Addr), bl.BlockTime, tx.Hash)
					conn.Do("ZADD", fmt.Sprintf("addr:%v:received", ntxo.Addr), bl.BlockTime, tx.Hash)

					conn.Do("HINCRBY", fmt.Sprintf("addr:%v:h", ntxo.Addr), "tr", ntxo.Value)

					txomut.Lock()
					txos = append(txos, ntxo)
					txomut.Unlock()

				}(bl, tx, pool, &total_tx_out, txo, txo_index)
			}

			//txis_cnt := uint32(0)
			// Skip the ins if it's a CoinBase Tx (1 TxIn for newly generated coins)
			if !(len(tx.TxIns) == 1 && tx.TxIns[0].InputVout == 0xffffffff) {

				for txi_index, txi := range tx.TxIns {
					txwg.Add(1)
					sem <- true
					go func(txi *blkparser.TxIn, bl *blkparser.Block, tx *blkparser.Tx, pool *redis.Pool, total_tx_in *uint64, txi_index int) {
						conn := pool.Get()
						defer conn.Close()
						defer func() {
							<-sem
						}()
						defer txwg.Done()

						ntxi := new(TxIn)
						ntxi.TxHash = tx.Hash
						ntxi.BlockHash = bl.Hash
						ntxi.BlockTime = bl.BlockTime
						ntxi.Index = uint32(txi_index)
						nprevout := new(PrevOut)
						nprevout.Vout = txi.InputVout
						nprevout.Hash = txi.InputHash
						ntxi.PrevOut = nprevout
						prevtxo := new(TxOutCached)

						prevtxocachedraw, err := ldb.Get(ro, []byte(fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))
						if err != nil {
							log.Printf("Err getting prevtxocached: %v", err)
						}

						if len(prevtxocachedraw) > 0 {
							if err := json.Unmarshal(prevtxocachedraw, prevtxo); err != nil {
								panic(err)
							}
						} else {
							// Shouldn't happen!
							//log.Println("Fallback to SSDB")
							prevtxoredisjson, err := redis.String(conn.Do("GET", fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))
							if err != nil {
								log.Printf("KEY:%v\n", fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout))
								panic(err)
							}
							prevtxoredis := new(TxOut)
							json.Unmarshal([]byte(prevtxoredisjson), prevtxoredis)

							prevtxo.Addr = prevtxoredis.Addr
							prevtxo.Value = prevtxoredis.Value
							//prevtxo.Id = prevtxomongo.Id.Hex()
						}

						ldb.Delete(wo, []byte(fmt.Sprintf("txo:%v:%v", txi.InputHash, txi.InputVout)))

						nprevout.Address = prevtxo.Addr
						nprevout.Value = prevtxo.Value

						txospent := new(TxoSpent)
						txospent.Spent = true
						txospent.BlockHeight = uint32(block_height)
						txospent.InputHash = tx.Hash
						txospent.InputIndex = uint32(txi_index)

						//total_tx_in+= uint(nprevout.Value)
						atomic.AddUint64(total_tx_in, nprevout.Value)

						tximut.Lock()
						txis = append(txis, ntxi)
						tximut.Unlock()
						//atomic.AddUint32(txis_cnt, 1)

						//log.Println("Starting update prev txo")
						ntxijson, _ := json.Marshal(ntxi)
						ntxikey := fmt.Sprintf("txi:%v:%v", tx.Hash, txi_index)

						txospentjson, _ := json.Marshal(txospent)

						conn.Do("SET", ntxikey, ntxijson)
						//conn.Send("ZADD", fmt.Sprintf("txi:%v", tx.Hash), txi_index, ntxikey)

						conn.Do("SET", fmt.Sprintf("txo:%v:%v:spent", txi.InputHash, txi.InputVout), txospentjson)

						conn.Do("ZADD", fmt.Sprintf("addr:%v", nprevout.Address), bl.BlockTime, tx.Hash)
						conn.Do("ZADD", fmt.Sprintf("addr:%v:sent", nprevout.Address), bl.BlockTime, tx.Hash)
						conn.Do("HINCRBY", fmt.Sprintf("addr:%v:h", nprevout.Address), "ts", nprevout.Value)
					}(txi, bl, tx, pool, &total_tx_in, txi_index)

				}
			}

			err := ldb.Write(wo, wb)
			if err != nil {
				log.Fatalf("Err write batch: %v", err)
			}
			wb.Clear()

			txwg.Wait()

			total_bl_out += total_tx_out

			ntx := new(Tx)
			ntx.Index = uint32(tx_index)
			ntx.Hash = tx.Hash
			ntx.Size = tx.Size
			ntx.LockTime = tx.LockTime
			ntx.Version = tx.Version
			ntx.TxInCnt = uint32(len(txis))
			ntx.TxOutCnt = uint32(len(txos))
			ntx.TotalOut = uint64(total_tx_out)
			ntx.TotalIn = uint64(total_tx_in)
			ntx.BlockHash = bl.Hash
			ntx.BlockHeight = block_height
			ntx.BlockTime = bl.BlockTime

			ntxjson, _ := json.Marshal(ntx)
			ntxjsonkey := fmt.Sprintf("tx:%v", ntx.Hash)
			conn.Do("SET", ntxjsonkey, ntxjson)
			conn.Do("ZADD", fmt.Sprintf("tx:%v:blocks", tx.Hash), bl.BlockTime, bl.Hash)
			conn.Do("ZADD", fmt.Sprintf("block:%v:txs", block.Hash), tx_index, ntxjsonkey)

			ntx.TxIns = txis
			ntx.TxOuts = txos
			txs = append(txs, ntx)
		}

		block.TotalBTC = uint64(total_bl_out)
		block.TxCnt = uint32(len(txs))

		blockjson, _ := json.Marshal(block)
		conn.Do("ZADD", "blocks", block.BlockTime, block.Hash)
		conn.Do("MSET", fmt.Sprintf("block:%v", block.Hash), blockjson, "height:latest", int(block_height), fmt.Sprintf("block:height:%v", block.Height), block.Hash)
		block.Txs = txs
		blockjsoncache, _ := json.Marshal(block)
		conn.Do("SET", fmt.Sprintf("block:%v:cached", block.Hash), blockjsoncache)

		if !running {
			log.Printf("Done. Stopped at height: %v.", block_height)
		}

		wg.Done()
	}
	wg.Wait()
}
