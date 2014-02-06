package btcplex

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
    "github.com/garyburd/redigo/redis"
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

// Fetch a block via bitcoind RPC API
func GetBlockRPC(conf *Config, block_height uint) (block *Block, txs []*Tx, err error) {
	// Get the block hash
	res, err := CallBitcoinRPC(conf.BitcoindRpcUrl, "getblockhash", 1, []interface{}{block_height})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	res, err = CallBitcoinRPC(conf.BitcoindRpcUrl, "getblock", 1, []interface{}{res["result"]})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	blockjson := res["result"].(map[string]interface{})

	block = new(Block)
	block.Hash = blockjson["hash"].(string)
	block.Height = block_height
	vertmp, _ := blockjson["version"].(json.Number).Int64()
	block.Version = uint32(vertmp)
	block.MerkleRoot = blockjson["merkleroot"].(string)
	block.Parent = blockjson["previousblockhash"].(string)
	sizetmp, _ := blockjson["size"].(json.Number).Int64()
	block.Size = uint32(sizetmp)
	noncetmp, _ := blockjson["nonce"].(json.Number).Int64()
	block.Nonce = uint32(noncetmp)
	btimetmp, _ := blockjson["time"].(json.Number).Int64()
	block.BlockTime = uint32(btimetmp)
	blockbits, _ := strconv.ParseInt(blockjson["bits"].(string), 16, 0)
	block.Bits = uint32(blockbits)
	block.TxCnt = uint32(len(blockjson["tx"].([]interface{})))
	fmt.Printf("Endblockrpc")
	txs = []*Tx{}
	tout := uint64(0)
	for _, txjson := range blockjson["tx"].([]interface{}) {
		tx, _ := GetTxRPC(conf, txjson.(string), block)
		tout += tx.TotalOut
		txs = append(txs, tx)
	}
	block.TotalBTC = uint64(tout * 1e8)
	return
}

// Fetch a transaction without additional info, used to fetch previous txouts when parsing txins
func GetRawTxRPC(conf *Config, tx_id string) (tx *Tx, err error) {
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
	vertmp, _ := txjson["version"].(json.Number).Int64()
	tx.Version = uint32(vertmp)
	ltimetmp, _ := txjson["locktime"].(json.Number).Int64()
	tx.LockTime = uint32(ltimetmp)
	tx.Size = uint32(len(txjson["hex"].(string)) / 2)
	//tx.

	total_tx_out := uint(0)
	total_tx_in := uint(0)

	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(TxIn)

			txinjsonprevout := new(PrevOut)
			txinjsonprevout.Hash = txijson.(map[string]interface{})["txid"].(string)
			vouttmp, _ := txijson.(map[string]interface{})["vout"].(json.Number).Int64()
			txinjsonprevout.Vout = uint32(vouttmp)
			txi.PrevOut = txinjsonprevout

			tx.TxIns = append(tx.TxIns, txi)
		}
	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(TxOut)
		valtmp, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
		txo.Value = uint64(valtmp * 1e8)
		if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nonstandard" {
			txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
		} else {
			txo.Addr = ""
		}
		tx.TxOuts = append(tx.TxOuts, txo)

		total_tx_out += uint(txo.Value)
	}

	tx.TxOutCnt = uint32(len(tx.TxOuts))
	tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TotalOut = uint64(total_tx_out)
	tx.TotalIn = uint64(total_tx_in)
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

			prevtx, _ := GetRawTxRPC(conf, txinjsonprevout.Hash)
			prevout := prevtx.TxOuts[txinjsonprevout.Vout]

			txinjsonprevout.Address = prevout.Addr
			txinjsonprevout.Value = prevout.Value

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
		txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
		tx.TxOuts = append(tx.TxOuts, txo)
		total_tx_out += uint64(txo.Value)
	}

	tx.TxOutCnt = uint32(len(tx.TxOuts))
	tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TotalOut = uint64(total_tx_out)
	tx.TotalIn = uint64(total_tx_in)
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

// Get unconfirmed transactions from memory pool, along with
// first seem time/block height, requires a recent bitcoind version
func GetUnconfirmedTxsRPC(conf *Config, pool *redis.Pool) (utxs []*Tx, err error) {    
	c := pool.Get()
	defer c.Close()
	unconfirmedtxsverbose, _ := GetRawMemPoolVerboseRPC(conf)
	unconfirmedtxs, _ := GetRawMemPoolRPC(conf)
	utxs = []*Tx{}
	for _, txid := range unconfirmedtxs {
		utxexists, _ := redis.Bool(c.Do("SADD", "btcplex:rawmempool", txid))
		if utxexists {
			tx, _ := GetTxRPC(conf, txid, &Block{})
			txmeta, txfound := unconfirmedtxsverbose[txid].(map[string]interface{})
			if !txfound {
				c.Do("SREM", "btcplex:rawmempool", txid)
				log.Println("Not found in verbose rawmempool")
			} else {
				fseentime, _ := txmeta["time"].(json.Number).Int64()
				tx.FirstSeenTime = uint32(fseentime)
				fseenheight, _ := txmeta["height"].(json.Number).Int64()
				tx.FirstSeenHeight = uint(fseenheight)
				utxs = append(utxs, tx)
				txjson, _ := json.Marshal(tx)
				c.Do("SET", fmt.Sprintf("btcplex:utx:%v", txid), string(txjson))
			}
		} else {
			txraw, _ := redis.String(c.Do("GET", fmt.Sprintf("btcplex:utx:%v", txid)))
			utx := new(Tx)
			if txraw != "" {
				err = json.Unmarshal([]byte(txraw), utx)
				if err != nil {
					return
				}
				utxs = append(utxs, utx)
			}
		}
		log.Println("Done")
		
	}
	return
}
