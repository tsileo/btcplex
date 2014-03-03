package btcplex

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "io/ioutil"
	"log"
	"net/http"
	"strconv"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
)

const GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"

// Helper to make call to bitcoind RPC API
func CallBitcoinRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(address,
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalf("ReadAll: %v", err)
	//	return nil, err
	//}
	result := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&result)
	//err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}
	return result, nil
}

func GetBlockHashRPC(conf *Config, height uint) string {
	// Get the block hash
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getblockhash", 1, []interface{}{height})
	if err != nil {
		return ""
	}
	return res["result"].(string)
}

func GetBlockCountRPC(conf *Config) uint {
	// Get the block hash
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getblockcount", 1, []interface{}{})
	if err != nil {
		return uint(0)
	}
	count, _ := res["result"].(json.Number).Int64()
	return uint(count)
}

func SaveBlockFromRPC(conf *Config, pool *redis.Pool, hash string) (block *Block, err error) {
	c := pool.Get()
	defer c.Close()
	var wg sync.WaitGroup
	sem := make(chan bool, 5)
	// Get the block hash
	//res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getblockhash", 1, []interface{}{block_height})
	//if err != nil {
	//	log.Fatalf("Err: %v", err)
	//}
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getblock", 1, []interface{}{hash})
	if err != nil {
		return
	}
	if res["result"] == nil {
		err = errors.New("Error fetching block")
		return
	}
	blockjson := res["result"].(map[string]interface{})

	block = new(Block)
	block.Hash = blockjson["hash"].(string)
	bheight, _ := blockjson["height"].(json.Number).Int64()
	block.Height = uint(bheight)
	block.Parent = blockjson["previousblockhash"].(string)

	prevheight := block.Height - 1
	prevhashtest := block.Parent
	prevnext := block.Hash
	for {
		prevkey := fmt.Sprintf("height:%v", prevheight)
		prevcnt, _ := redis.Int(c.Do("ZCARD", prevkey))
		// SSDB doesn't support negative slice yet
		prevs, _ := redis.Strings(c.Do("ZRANGE", prevkey, 0, prevcnt-1))
		if len(prevs) == 0 {
			break
		}
		for _, cprevhash := range prevs {
			if cprevhash == prevhashtest {
				// current block parent
				prevhashtest, _ = redis.String(c.Do("HGET", fmt.Sprintf("block:%v:h", cprevhash), "parent"))
				// Set main to 1 and the next => prevnext
				c.Do("HMSET", fmt.Sprintf("block:%v:h", cprevhash), "main", true, "next", prevnext)
				c.Do("SET", fmt.Sprintf("block:height:%v", prevheight), cprevhash)
				prevnext = cprevhash
			} else {
				// Set main to 0
				c.Do("HSET", fmt.Sprintf("block:%v:h", cprevhash), "main", false)
				oblock, _ := GetBlockCachedByHash(pool, cprevhash)
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

	vertmp, _ := blockjson["version"].(json.Number).Int64()
	block.Version = uint32(vertmp)
	block.MerkleRoot = blockjson["merkleroot"].(string)
	sizetmp, _ := blockjson["size"].(json.Number).Int64()
	block.Size = uint32(sizetmp)
	noncetmp, _ := blockjson["nonce"].(json.Number).Int64()
	block.Nonce = uint32(noncetmp)
	btimetmp, _ := blockjson["time"].(json.Number).Int64()
	block.BlockTime = uint32(btimetmp)
	blockbits, _ := strconv.ParseInt(blockjson["bits"].(string), 16, 0)
	block.Bits = uint32(blockbits)
	block.TxCnt = uint32(len(blockjson["tx"].([]interface{})))
	tout := uint64(0)
	txs := []*Tx{}
	var txmut sync.Mutex
	for txindex, txjson := range blockjson["tx"].([]interface{}) {
		sem <- true
		wg.Add(1)
		go func(txjson interface{}, tout *uint64, block *Block, txs *[]*Tx) {
			defer wg.Done()
			defer func() { <-sem }()
			tx, _ := SaveTxFromRPC(conf, pool, txjson.(string), block, txindex)
			//(conf *Config, pool *redis.Pool, tx_id string, block *Block, tx_index int)
			atomic.AddUint64(tout, tx.TotalOut)
			txmut.Lock()
			*txs = append(*txs, tx)
			txmut.Unlock()
		}(txjson, &tout, block, &txs)
	}
	wg.Wait()
	block.TotalBTC = uint64(tout)

	c.Do("ZADD", fmt.Sprintf("height:%v", block.Height), block.BlockTime, block.Hash)
	c.Do("HMSET", fmt.Sprintf("block:%v:h", block.Hash), "parent", block.Parent, "height", block.Height, "main", true)
	blockjson2, _ := json.Marshal(block)
	c.Do("ZADD", "blocks", block.BlockTime, block.Hash)
	c.Do("MSET", fmt.Sprintf("block:%v", block.Hash), blockjson2, "height:latest", int(block.Height), fmt.Sprintf("block:height:%v", block.Height), block.Hash)
	By(TxIndex).Sort(txs)
	block.Txs = txs
	fullblockjson, _ := json.Marshal(block)
	c.Do("SET", fmt.Sprintf("block:%v:cached", block.Hash), fullblockjson)
	return
}

// Fetch a transaction without additional info, used to fetch previous txouts when parsing txins
func GetTxOutRPC(conf *Config, tx_id string, txo_vout uint32) (txo *TxOut, err error) {
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	// Get the TX from bitcoind RPC API
	res_tx, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getrawtransaction", 1, []interface{}{tx_id, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	txjson := res_tx["result"].(map[string]interface{})

	txojson := txjson["vout"].([]interface{})[txo_vout]
	txo = new(TxOut)
	valtmp, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
	txo.Value = FloatToUint(valtmp)
	if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nonstandard" {
		txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
		if txoisinterface {
			txo.Addr = txodata[0].(string)
		} else {
			txo.Addr = ""
		}
	} else {
		txo.Addr = ""
	}
	txospent := new(TxoSpent)
	txospent.Spent = false
	txo.Spent = txospent
	return
}

// Fetch a transaction via bticoind RPC API
func GetTxRPC(conf *Config, tx_id string, block *Block) (tx *Tx, err error) {
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	// Get the TX from bitcoind RPC API
	res_tx, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getrawtransaction", 1, []interface{}{tx_id, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	txjson := res_tx["result"].(map[string]interface{})

	tx = new(Tx)
	tx.Hash = tx_id
	tx.BlockTime = block.BlockTime
	tx.BlockHeight = block.Height
	tx.BlockHash = block.Hash
	vertmp, _ := txjson["version"].(json.Number).Int64()
	tx.Version = uint32(vertmp)
	ltimetmp, _ := txjson["locktime"].(json.Number).Int64()
	tx.LockTime = uint32(ltimetmp)
	tx.Size = uint32(len(txjson["hex"].(string)) / 2)

	total_tx_out := uint64(0)
	total_tx_in := uint64(0)

	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(TxIn)
			txinjsonprevout := new(PrevOut)
			txinjsonprevout.Hash = txijson.(map[string]interface{})["txid"].(string)
			tmpvout, _ := txijson.(map[string]interface{})["vout"].(json.Number).Int64()
			txinjsonprevout.Vout = uint32(tmpvout)

			// Check if bitcoind is patched to fetch value/address without additional RPC call
			// cf. README
			_, bitcoindPatched := txijson.(map[string]interface{})["value"]
			if bitcoindPatched {
				pval, _ := txijson.(map[string]interface{})["value"].(json.Number).Float64()
				txinjsonprevout.Address = txijson.(map[string]interface{})["address"].(string)
				txinjsonprevout.Value = FloatToUint(pval)
			} else {
				prevout, _ := GetTxOutRPC(conf, txinjsonprevout.Hash, txinjsonprevout.Vout)
				txinjsonprevout.Address = prevout.Addr
				txinjsonprevout.Value = prevout.Value
			}

			total_tx_in += uint64(txinjsonprevout.Value)

			txi.PrevOut = txinjsonprevout

			tx.TxIns = append(tx.TxIns, txi)

			// TODO handle txi from this TX
		}
	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(TxOut)
		txoval, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
		txo.Value = uint64(txoval * 1e8)
		//txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
		if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nonstandard" {
			txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
			if txoisinterface {
				txo.Addr = txodata[0].(string)
			} else {
				txo.Addr = ""
			}
		} else {
			txo.Addr = ""
		}
		tx.TxOuts = append(tx.TxOuts, txo)
		txospent := new(TxoSpent)
		txospent.Spent = false
		txo.Spent = txospent
		total_tx_out += uint64(txo.Value)
	}

	tx.TxOutCnt = uint32(len(tx.TxOuts))
	tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TotalOut = uint64(total_tx_out)
	tx.TotalIn = uint64(total_tx_in)
	return
}

// Fetch a transaction via bticoind RPC API
func SaveTxFromRPC(conf *Config, pool *redis.Pool, tx_id string, block *Block, tx_index int) (tx *Tx, err error) {
	c := pool.Get()
	defer c.Close()
	var wg sync.WaitGroup
	var tximut, txomut sync.Mutex
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	// Get the TX from bitcoind RPC API
	res_tx, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getrawtransaction", 1, []interface{}{tx_id, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	txjson := res_tx["result"].(map[string]interface{})

	tx = new(Tx)
	tx.Index = uint32(tx_index)
	tx.Hash = tx_id
	tx.BlockTime = block.BlockTime
	tx.BlockHeight = block.Height
	tx.BlockHash = block.Hash
	vertmp, _ := txjson["version"].(json.Number).Int64()
	tx.Version = uint32(vertmp)
	ltimetmp, _ := txjson["locktime"].(json.Number).Int64()
	tx.LockTime = uint32(ltimetmp)
	tx.Size = uint32(len(txjson["hex"].(string)) / 2)

	total_tx_out := uint64(0)
	total_tx_in := uint64(0)

	sem := make(chan bool, 50)
	for txiindex, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			wg.Add(1)
			sem <- true
			go func(pool *redis.Pool, txijson interface{}, txiindex int, total_tx_in *uint64, tx *Tx, block *Block) {
				defer wg.Done()
				defer func() { <-sem }()
				c := pool.Get()
				defer c.Close()
				txi := new(TxIn)
				txinjsonprevout := new(PrevOut)
				txinjsonprevout.Hash = txijson.(map[string]interface{})["txid"].(string)
				tmpvout, _ := txijson.(map[string]interface{})["vout"].(json.Number).Int64()
				txinjsonprevout.Vout = uint32(tmpvout)

				// Check if bitcoind is patched to fetch value/address without additional RPC call
				// cf. README
				_, bitcoindPatched := txijson.(map[string]interface{})["value"]
				if bitcoindPatched {
					pval, _ := txijson.(map[string]interface{})["value"].(json.Number).Float64()
					txinjsonprevout.Address = txijson.(map[string]interface{})["address"].(string)
					txinjsonprevout.Value = FloatToUint(pval)
				} else {
					prevout, _ := GetTxOutRPC(conf, txinjsonprevout.Hash, txinjsonprevout.Vout)

					txinjsonprevout.Address = prevout.Addr
					txinjsonprevout.Value = prevout.Value
				}
				atomic.AddUint64(total_tx_in, uint64(txinjsonprevout.Value))

				txi.PrevOut = txinjsonprevout

				tximut.Lock()
				tx.TxIns = append(tx.TxIns, txi)
				tximut.Unlock()

				txospent := new(TxoSpent)
				txospent.Spent = true
				txospent.BlockHeight = uint32(block.Height)
				txospent.InputHash = tx.Hash
				txospent.InputIndex = uint32(txiindex)

				ntxijson, _ := json.Marshal(txi)
				ntxikey := fmt.Sprintf("txi:%v:%v", tx.Hash, txiindex)

				txospentjson, _ := json.Marshal(txospent)

				c.Do("SET", ntxikey, ntxijson)
				//conn.Send("ZADD", fmt.Sprintf("txi:%v", tx.Hash), txi_index, ntxikey)

				c.Do("SET", fmt.Sprintf("txo:%v:%v:spent", txinjsonprevout.Hash, txinjsonprevout.Vout), txospentjson)

				c.Do("ZADD", fmt.Sprintf("addr:%v", txinjsonprevout.Address), block.BlockTime, tx.Hash)
				c.Do("ZADD", fmt.Sprintf("addr:%v:sent", txinjsonprevout.Address), block.BlockTime, tx.Hash)
				c.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txinjsonprevout.Address), "ts", txinjsonprevout.Value)

			}(pool, txijson, txiindex, &total_tx_in, tx, block)
		}
	}
	for txo_index, txojson := range txjson["vout"].([]interface{}) {
		wg.Add(1)
		sem <- true
		go func(pool *redis.Pool, txojson interface{}, txo_index int, total_tx_out *uint64, tx *Tx, block *Block) {
			defer wg.Done()
			defer func() { <-sem }()
			c := pool.Get()
			defer c.Close()
			txo := new(TxOut)
			txoval, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
			txo.Value = FloatToUint(txoval)
			//txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)

			if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nonstandard" {
				txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
				if txoisinterface {
					txo.Addr = txodata[0].(string)
				} else {
					txo.Addr = ""
				}
			} else {
				txo.Addr = ""
			}

			txomut.Lock()
			tx.TxOuts = append(tx.TxOuts, txo)
			txomut.Unlock()
			txospent := new(TxoSpent)
			txospent.Spent = false
			txo.Spent = txospent
			//total_tx_out += uint64(txo.Value)
			atomic.AddUint64(total_tx_out, uint64(txo.Value))

			ntxojson, _ := json.Marshal(txo)
			ntxokey := fmt.Sprintf("txo:%v:%v", tx.Hash, txo_index)
			c.Do("SET", ntxokey, ntxojson)
			//conn.Send("ZADD", fmt.Sprintf("txo:%v", tx.Hash), txo_index, ntxokey)
			c.Do("ZADD", fmt.Sprintf("addr:%v", txo.Addr), block.BlockTime, tx.Hash)
			c.Do("ZADD", fmt.Sprintf("addr:%v:received", txo.Addr), block.BlockTime, tx.Hash)
			c.Do("HINCRBY", fmt.Sprintf("addr:%v:h", txo.Addr), "tr", txo.Value)
		}(pool, txojson, txo_index, &total_tx_out, tx, block)

	}

	wg.Wait()

	tx.TxOutCnt = uint32(len(tx.TxOuts))
	tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TotalOut = uint64(total_tx_out)
	tx.TotalIn = uint64(total_tx_in)

	ntxjson, _ := json.Marshal(tx)
	ntxjsonkey := fmt.Sprintf("tx:%v", tx.Hash)
	c.Do("SET", ntxjsonkey, ntxjson)
	c.Do("ZADD", fmt.Sprintf("block:%v:txs", block.Hash), tx_index, ntxjsonkey)
	c.Do("ZADD", fmt.Sprintf("tx:%v:blocks", tx.Hash), tx.BlockTime, block.Hash)
	return
}

func GetRawMemPoolRPC(conf *Config) (unconfirmedtxs []string, err error) {
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getrawmempool", 1, []interface{}{})
	if err != nil {
		return
	}
	unconfirmedtxs = []string{}
	for _, txid := range res["result"].([]interface{}) {
		unconfirmedtxs = append(unconfirmedtxs, txid.(string))
	}
	return
}

func GetRawMemPoolVerboseRPC(conf *Config) (unconfirmedtxs map[string]interface{}, err error) {
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getrawmempool", 1, []interface{}{true})
	if err != nil {
		return
	}
	unconfirmedtxs = res["result"].(map[string]interface{})
	return
}
