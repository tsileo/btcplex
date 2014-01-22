package btcplex

import (
	"encoding/json"
	"fmt"
	"github.com/jmhodges/levigo"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
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
	block.Version = uint32(blockjson["version"].(float64))
	block.MerkleRoot = blockjson["merkleroot"].(string)
	block.Parent = blockjson["previousblockhash"].(string)
	block.Size = uint32(blockjson["size"].(float64))
	block.Nonce = uint32(blockjson["nonce"].(float64))
	block.BlockTime = uint32(blockjson["time"].(float64))
	blockbits, _ := strconv.ParseInt(blockjson["bits"].(string), 16, 0)
	block.Bits = uint32(blockbits)
	block.TxCnt = uint32(len(blockjson["tx"].([]interface{})))
	fmt.Printf("Endblockrpc")
	txs = []*Tx{}
	tout := float64(0)
	for _, txjson := range blockjson["tx"].([]interface{}) {
		tx, itout, _ := GetTxRPC(conf, txjson.(string), block)
		tout += itout
		txs = append(txs, tx)
	}
	totalstr := fmt.Sprintf("%.8f", tout)
	totalbtc, _ := strconv.ParseFloat(totalstr, 0)
	block.TotalBTC = uint64(totalbtc * 1e8)
	return
}

// Fetch a transaction without additional info, used to fetch previous txouts when parsing txins
func QuickTxRPC(conf *Config, tx_id string) (tx *Tx, err error) {
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
	tx.Version = uint32(txjson["version"].(float64))
	tx.LockTime = uint32(txjson["locktime"].(float64))
	tx.Size = uint32(len(txjson["hex"].(string)) / 2)
	//tx.

	total_tx_out := uint(0)
	total_tx_in := uint(0)

	ro := levigo.NewReadOptions()
	defer ro.Close()

	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(TxIn)

			txinjsonprevout := new(PrevOut)
			txinjsonprevout.Hash = txijson.(map[string]interface{})["txid"].(string)
			txinjsonprevout.Vout = uint32(txijson.(map[string]interface{})["vout"].(float64))
			txi.PrevOut = txinjsonprevout

			tx.TxIns = append(tx.TxIns, txi)
		}
	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(TxOut)
		txo.Value = uint64(txojson.(map[string]interface{})["value"].(float64) * 1e8)
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
func GetTxRPC(conf *Config, tx_id string, block *Block) (tx *Tx, tout float64, err error) {
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
	tx.Version = uint32(txjson["version"].(float64))
	tx.LockTime = uint32(txjson["locktime"].(float64))
	tx.Size = uint32(len(txjson["hex"].(string)) / 2)

	total_tx_out := uint(0)
	total_tx_in := uint(0)
	tout = float64(0)
	ro := levigo.NewReadOptions()
	defer ro.Close()

	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(TxIn)
			txinjsonprevout := new(PrevOut)
			txinjsonprevout.Hash = txijson.(map[string]interface{})["txid"].(string)
			txinjsonprevout.Vout = uint32(txijson.(map[string]interface{})["vout"].(float64))

			prevtx, _ := QuickTxRPC(conf, txinjsonprevout.Hash)
			prevout := prevtx.TxOuts[txinjsonprevout.Vout]

			txinjsonprevout.Address = prevout.Addr
			txinjsonprevout.Value = prevout.Value

			total_tx_in += uint(txinjsonprevout.Value)

			txi.PrevOut = txinjsonprevout

			tx.TxIns = append(tx.TxIns, txi)

			// TODO handle txi from this TX
		}
	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(TxOut)
		txo.Value = uint64(txojson.(map[string]interface{})["value"].(float64) * 1e8)
		txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
		tx.TxOuts = append(tx.TxOuts, txo)

		total_tx_out += uint(txo.Value)
		tout += txojson.(map[string]interface{})["value"].(float64)
	}

	tx.TxOutCnt = uint32(len(tx.TxOuts))
	tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TotalOut = uint64(total_tx_out)
	tx.TotalIn = uint64(total_tx_in)
	return
}
