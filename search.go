package btcplex

import (
    "github.com/garyburd/redigo/redis"
	"strconv"
	"fmt"
)

func IsBlockHeight(rpool *redis.Pool, q string) (s bool, res string) {
	height, err := strconv.ParseUint(q, 10, 0)
	if err != nil {
		return false, ""
	}
	fmt.Sprintf("%v", height)
	//block, err := GetBlockByHeight(db, uint(height))
	//if err != nil {
	//	return false, ""
	//}
	//return true, block.Hash
	return false, ""
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
		return false, q
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
