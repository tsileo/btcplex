package btcplex

import (
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/pmylund/go-cache"
	"log"
	"strconv"
	"strings"
)

type AddressData struct {
	Address       string `json:"address"`
	TxCnt         uint64 `json:"n_tx"`
	ReceivedCnt   uint64 `json:"-"`
	SentCnt       uint64 `json:"-"`
	TotalReceived uint64 `json:"total_received"`
	TotalSent     uint64 `json:"total_sent"`
	FinalBalance  uint64 `json:"final_balance"`
	Txs           []*Tx  `json:"txs"`
}

// Timestamp of the block an address was first confirmed in
func AddressFirstSeen(db *levigo.DB, addr string) (timestamp uint, err error) {
	timestamp = uint(0)
	ro := levigo.NewReadOptions()
	defer ro.Close()
	it := db.NewIterator(ro)
	defer it.Close()
	it.Seek([]byte(fmt.Sprintf("%s-txo!", addr)))
	if it.Valid() {
		k := string(it.Key()[:])
		data := strings.Split(k, "!")
		r, _ := db.Get(ro, []byte(fmt.Sprintf("tx!%s", data[1])))
		data = strings.Split(string(r[:]), ":")
		ts, _ := strconv.Atoi(data[0])
		timestamp = uint(ts)
		return
	} else {
		err = errors.New("Address not found")
		return
	}
}

// Get the total number of bitcoins received by an address (in satoshi)
func GetReceivedByAddress(db *levigo.DB, addr string) (total uint, err error) {
	total = uint(0)
	ro := levigo.NewReadOptions()
	defer ro.Close()
	start := []byte(fmt.Sprintf("%s-txo!", addr))
	end := []byte(fmt.Sprintf("%s-txo!\xff", addr))
	txos, _ := GetRange(db, start, end)
	for _, txo := range txos {
		val, _ := strconv.Atoi(txo.Value)
		total += uint(val)
	}
	return
}

func GetSentByAddress(db *levigo.DB, addr string) (total uint, err error) {
	total = uint(0)
	ro := levigo.NewReadOptions()
	defer ro.Close()
	start := []byte(fmt.Sprintf("%s-txo!", addr))
	end := []byte(fmt.Sprintf("%s-txo!\xff", addr))
	txos, _ := GetRange(db, start, end)
	for _, txo := range txos {
		log.Println(txo.Key)
		new_key := strings.Replace(txo.Key, "-txo", "-txo-spent", 1)
		log.Println(new_key)
		r, _ := db.Get(ro, []byte(new_key))
		log.Println(string(r[:]))
		val, _ := strconv.Atoi(txo.Value)
		total += uint(val)
	}
	return
}

func GetAddressCached(c *cache.Cache, db *levigo.DB, address string) (addressdata *AddressData, err error) {
	cachekey := fmt.Sprintf("address%v", address)
	cached, found := c.Get(cachekey)
	if found {
		return cached.(*AddressData), nil
	} else {
		addressdata, err := GetAddress(db, address)
		c.Set(cachekey, addressdata, 0)
		return addressdata, err
	}
}

// Return address summary and history
// TODO sort txs
func GetAddress(db *levigo.DB, address string) (addressdata *AddressData, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	txos, _ := GetRange(db, []byte(fmt.Sprintf("%v-txo!", address)), []byte(fmt.Sprintf("%v-txo!\xff", address)))

	addressdata = new(AddressData)
	txsindex := map[string]interface{}{}

	txreceived := map[string]interface{}{}
	txsent := map[string]interface{}{}

	txs := []*Tx{}

	for _, txo := range txos {
		//txo_value, _ := strconv.ParseUint(txo.Value, 10, 0)
		txdata := strings.Split(txo.Key, "!")

		_, inindex := txsindex[txdata[1]]
		if !inindex {

			tx, _ := GetTx(db, txdata[1])
			fmt.Printf("GETTX")
			//txaddressinfo := new(TxAddressInfo)
			//txaddressinfo.InTxOut = true
			//txaddressinfo.InTxIn = false
			//txaddressinfo.Value = uint64(txo_value)
			//tx.TxAddressInfo = txaddressinfo

			// =>
			txsindex[txdata[1]] = tx

			rtxo_index, _ := strconv.ParseUint(txdata[2], 10, 0)
			//txospent, txoerr := GetTxoSpent(db, address, tx.Hash, txo_index)
			txospent := tx.TxOuts[rtxo_index].Spent
			if txospent.Spent {
				_, inindex2 := txsindex[txospent.InputHash]
				if !inindex2 {
					ntx, _ := GetTx(db, txospent.InputHash)
					fmt.Printf("GETXI")
					txsindex[txospent.InputHash] = ntx
				}
			}
		}

	}

	//"n_tx":24,
	//"total_received":515863600,
	//"total_sent":515863600,
	//"final_balance":0,

	totalreceived := uint64(0)
	totalsent := uint64(0)

	for _, tx := range txsindex {
		ctx := tx.(*Tx)
		txaddressinfo := new(TxAddressInfo)
		for _, txi := range ctx.TxIns {
			if txi.PrevOut.Address == address {
				txaddressinfo.InTxIn = true
				txaddressinfo.Value -= int64(txi.PrevOut.Value)
				totalsent += txi.PrevOut.Value
				_, inindex := txsent[ctx.Hash]
				if !inindex {
					txsent[ctx.Hash] = true
				}
			}
		}
		for _, txo := range ctx.TxOuts {
			if txo.Addr == address {
				txaddressinfo.InTxOut = true
				txaddressinfo.Value += int64(txo.Value)
				totalreceived += txo.Value
				_, inindex := txreceived[ctx.Hash]
				if !inindex {
					txreceived[ctx.Hash] = true
				}
			}
		}
		ctx.TxAddressInfo = txaddressinfo
		txs = append(txs, ctx)
	}

	finalbalance := totalreceived - totalsent

	addressdata.Txs = txs
	addressdata.FinalBalance = finalbalance
	addressdata.TotalSent = totalsent
	addressdata.TotalReceived = totalreceived
	addressdata.TxCnt = uint64(len(txs))
	addressdata.Address = address
	addressdata.SentCnt = uint64(len(txsent))
	addressdata.ReceivedCnt = uint64(len(txreceived))

	return
}
