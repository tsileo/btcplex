package btcplex

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

// Use SSDB
func IsBlockHeight(rpool *redis.Pool, q string) (s bool, res string) {
	height, err := strconv.ParseUint(q, 10, 0)
	if err != nil {
		return false, ""
	}
	c := rpool.Get()
	defer c.Close()
	hash, err := redis.String(c.Do("GET", fmt.Sprintf("block:height:%v", height)))
	if err != nil {
		return false, ""
	}
	return true, hash
}

func IsBlockHash(rpool *redis.Pool, q string) (s bool, res string) {
	if len(q) != 64 {
		return false, ""
	}
	block, err := GetBlockByHash(rpool, q)
	if err != nil {
		return false, ""
	}
	return true, block.Hash
}

func IsTxHash(rpool *redis.Pool, q string) (s bool, res string) {
	if len(q) != 64 {
		return false, ""
	}
	tx, err := GetTx(rpool, q)
	if err != nil {
		return false, ""
	}
	return true, tx.Hash
}

// Check if the string is a valid Bitcoin address
func IsAddress(q string) (s bool, res string) {
	valid, _ := ValidA58([]byte(q))
	if valid {
		return true, q
	}
	return false, ""
}

// Check if the Tx is in Redis (not SSDB, Redis!) (in rawmempool)
func IsUnconfirmedTx(pool *redis.Pool, hash string) (status bool, res string) {
	if len(hash) != 64 {
		return false, ""
	}
	c := pool.Get()
	defer c.Close()
	txkey := fmt.Sprintf("btcplex:utx:%v", hash)
	status, _ = redis.Bool(c.Do("EXISTS", txkey))
	if status {
		return status, hash
	}
	return status, ""
}
