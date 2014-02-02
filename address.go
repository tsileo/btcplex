package btcplex

import (
    "fmt"
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

type AddressHash struct {
    TotalSent       int `redis:"ts"`
    TotalReceived   int `redis:"tr"`
}

func GetAddress(rpool *redis.Pool, address string) (addressdata *AddressData, err error) {
    c := rpool.Get()
    defer c.Close()

    addressdata = new(AddressData)
    txs1 := []*Tx{}
    txs := []*Tx{}

    zkey := fmt.Sprintf("addr:%v", address)
    txscnt, _ := redis.Int(c.Do("ZCARD", zkey))
    data, _ := redis.Strings(c.Do("ZREVRANGE", zkey, 0, 50))

    addressh := new(AddressHash)
    v, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("addr:%v:h", address)))
    if err != nil {
        panic(err)
    }
    if err := redis.ScanStruct(v, addressh); err != nil {
        panic(err)
    }

    for _, txd := range data {
        tx, _ := GetTx(rpool, txd)
        txs1 = append(txs1, tx)
    }

    totalreceived := uint64(addressh.TotalReceived)
    totalsent := uint64(addressh.TotalSent)
    finalbalance := totalreceived - totalreceived
    sentcnt, _ := redis.Int(c.Do("ZCARD", fmt.Sprintf("addr:%v:sent", address)))
    receivedcnt, _ := redis.Int(c.Do("ZCARD", fmt.Sprintf("addr:%v:received", address)))

    for _, ctx := range txs1 {
        txaddressinfo := new(TxAddressInfo)
        for _, txi := range ctx.TxIns {
            if txi.PrevOut.Address == address {
                txaddressinfo.InTxIn = true
                txaddressinfo.Value -= int64(txi.PrevOut.Value)
            }
        }
        for _, txo := range ctx.TxOuts {
            if txo.Addr == address {
                txaddressinfo.InTxOut = true
                txaddressinfo.Value += int64(txo.Value)
            }
        }
        ctx.TxAddressInfo = txaddressinfo
        txs = append(txs, ctx)
    }

    //By(TxBlockTime).Sort(txs)

    addressdata.Txs = txs
    addressdata.FinalBalance = finalbalance
    addressdata.TotalSent = totalsent
    addressdata.TotalReceived = totalreceived
    addressdata.TxCnt = uint64(txscnt)
    addressdata.Address = address
    addressdata.SentCnt = uint64(sentcnt)
    addressdata.ReceivedCnt = uint64(receivedcnt)

    return
}

// Return the block time at which the address first appeared
func AddressFirstSeen(rpool *redis.Pool, address string) (firstseen uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    zkey := fmt.Sprintf("addr:%v", address)
    data, _ := redis.Strings(c.Do("ZRANGE", zkey, 0, 0, "withscores"))
    txoutblocktime, _ := strconv.Atoi(data[1])
    firstseen = uint64(txoutblocktime)
    return
}

func GetReceivedByAddress(rpool *redis.Pool, address string) (total uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    addressh := new(AddressHash)
    v, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("addr:%v:h", address)))
    if err != nil {
        panic(err)
    }
    if err := redis.ScanStruct(v, addressh); err != nil {
        panic(err)
    }
    total = uint64(addressh.TotalReceived)
    return
}

func GetSentByAddress(rpool *redis.Pool, address string) (total uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    addressh := new(AddressHash)
    v, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("addr:%v:h", address)))
    if err != nil {
        panic(err)
    }
    if err := redis.ScanStruct(v, addressh); err != nil {
        panic(err)
    }
    total = uint64(addressh.TotalSent)
    return
}

func AddressBalance(rpool *redis.Pool, address string) (balance uint64, err error) {
    c := rpool.Get()
    defer c.Close()

    addressh := new(AddressHash)
    v, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("addr:%v:h", address)))
    if err != nil {
        panic(err)
    }
    if err := redis.ScanStruct(v, addressh); err != nil {
        panic(err)
    }
    balance = uint64(addressh.TotalReceived - addressh.TotalSent)
    return
}
