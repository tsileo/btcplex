package btcplex

import (
    "encoding/json"
    "fmt"
    "strings"
    "strconv"
    "github.com/garyburd/redigo/redis"
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

func GetAddress(rpool *redis.Pool, address string) (addressdata *AddressData, err error) {
    c := rpool.Get()
    defer c.Close()

    addressdata = new(AddressData)
    txsindex := map[string]interface{}{}

    txreceived := map[string]interface{}{}
    txsent := map[string]interface{}{}

    txs := []*Tx{}

    zkey := fmt.Sprintf("addr:%v", address)
    cnt, _ := redis.Int(c.Do("ZCARD", zkey))
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, cnt))

    for _, txd := range data {
        txds := strings.Split(txd, ":")
        chash := txds[1]
        _, inindex := txsindex[chash]
        if !inindex {
            tx, _ := GetTx(rpool, txds[1])
            txsindex[chash] = tx
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

    //By(TxBlockTime).Sort(txs)

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
func AddressFirstSeen(rpool *redis.Pool, address string) (firstseen uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    zkey := fmt.Sprintf("addr:%v", address)
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, 0, "withscores"))
    txoutblocktime, _ := strconv.Atoi(data[1])
    firstseen = uint64(txoutblocktime)
    return
}

func GetReceivedByAddress(rpool *redis.Pool, address string) (total uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    zkey := fmt.Sprintf("addr:%v", address)
    cnt, _ := redis.Int(c.Do("ZCARD", zkey))
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, cnt))

    for _, txd := range data {
        if strings.Contains(txd, "txo:") {
            ctxo := new(TxOut)
            txoutjson, _ := redis.String(c.Do("GET", txd))
            err = json.Unmarshal([]byte(txoutjson), ctxo)
            total += uint64(ctxo.Value)
        }
    }

    return
}

func GetSentByAddress(rpool *redis.Pool, address string) (total uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    zkey := fmt.Sprintf("addr:%v", address)
    cnt, _ := redis.Int(c.Do("ZCARD", zkey))
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, cnt))

    for _, txd := range data {
        if strings.Contains(txd, "txi:") {
            ctxi := new(TxIn)
            txinjson, _ := redis.String(c.Do("GET", txd))
            err = json.Unmarshal([]byte(txinjson), ctxi)
            total += uint64(ctxi.PrevOut.Value)
        }
    }

    return
}

func AddressBalance(rpool *redis.Pool, address string) (balance uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    zkey := fmt.Sprintf("addr:%v", address)
    cnt, _ := redis.Int(c.Do("ZCARD", zkey))
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, cnt))

    for _, txd := range data {
        if strings.Contains(txd, "txo:") {
            ctxo := new(TxOut)
            txoutjson, _ := redis.String(c.Do("GET", txd))
            err = json.Unmarshal([]byte(txoutjson), ctxo)
            balance += uint64(ctxo.Value)
        }
        if strings.Contains(txd, "txi:") {
            ctxi := new(TxIn)
            txinjson, _ := redis.String(c.Do("GET", txd))
            err = json.Unmarshal([]byte(txinjson), ctxi)
            balance -= uint64(ctxi.PrevOut.Value)
        }
    }

    return
}
