package btcplex

import (
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

type AddressData struct {
        Address string `json:"address"`
        TxCnt uint64 `json:"n_tx"`
        ReceivedCnt uint64 `json:"-"`
        SentCnt uint64 `json:"-"`
        TotalReceived uint64 `json:"total_received"`
        TotalSent uint64 `json:"total_sent"`
        FinalBalance uint64 `json:"final_balance"`
        Txs []*Tx `json:"txs"`
}

func GetAddress(db *mgo.Database, address string) (addressdata *AddressData, err error) {
        txos := []*TxOut{}
        err = db.C("txos").Find(bson.M{"addr": address}).All(&txos)
        if err != nil {
        	panic(err)
        }

        addressdata = new(AddressData)
        txsindex := map[string]interface{}{}

        txreceived := map[string]interface{}{}
        txsent := map[string]interface{}{}

        txs := []*Tx{}

        for _, txo := range txos {
        		_, inindex := txsindex[txo.TxHash]
                if !inindex {
                        tx, _ := GetTx(db, txo.TxHash)
                        tx.Build(db)
                        txsindex[txo.TxHash] = tx
						if txo.Spent.Spent {
                                _, inindex2 := txsindex[txo.Spent.InputHash]
                                if !inindex2 {
                                        ntx, _ := GetTx(db, txo.Spent.InputHash)
                                        ntx.Build(db)
                                        txsindex[txo.Spent.InputHash] = ntx
                                }
                        }
                }

        }

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

        By(TxBlockTime).Sort(txs)

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

// Return the block time at which the address first appeared  
func AddressFirstSeen(db *mgo.Database, addr string) (firstseen uint64, err error) {
    txout := new(TxOut)
	err = db.C("txos").Find(bson.M{"addr": addr}).Sort("blocktime").Limit(1).One(txout)
	if err != nil {
		return
	}
	firstseen = uint64(txout.BlockTime)
	return 
}

func GetReceivedByAddress(db *mgo.Database, addr string) (total uint64, err error) {
	txos := []*TxOut{}
	err = db.C("txos").Find(bson.M{"addr": addr}).All(&txos)
	if err != nil {
		return
	}
	for _, txo := range txos {
		total +=uint64(txo.Value)
	}
	return
}

func GetSentByAddress(db *mgo.Database, addr string) (total uint64, err error) {
	txos := []*TxOut{}
	err = db.C("txos").Find(bson.M{"addr": addr}).All(&txos)
	if err != nil {
		return
	}
	for _, txo := range txos {
		if txo.Spent.Spent {
			total +=uint64(txo.Value)	
		}
	}
	return
}

func AddressBalance(db *mgo.Database, addr string) (balance uint64, err error) {
	txos := []*TxOut{}
	err = db.C("txos").Find(bson.M{"addr": addr}).All(&txos)
	if err != nil {
		return
	}
	for _, txo := range txos {
		if !txo.Spent.Spent {
			balance +=uint64(txo.Value)	
		}
	}
	return

}
