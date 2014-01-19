package btcplex

import (
	"labix.org/v2/mgo"
	"strconv"
)

func IsBlockHeight(db *mgo.Database, q string) (s bool, res string) {
	height, err := strconv.ParseUint(q, 10, 0)
	if err != nil {
		return false, ""
	}
	block, err := GetBlockByHeight(db, uint(height))
	if err != nil {
		return false, ""
	}
	return true, block.Hash
}

func IsBlockHash(db *mgo.Database, q string) (s bool, res string) {
	if len(q) != 64 {
		return false, ""
	}
	block, err := GetBlockByHash(db, q)
	if err != nil {
		return false, ""
	}
	return true, block.Hash
}

func IsTxHash(db *mgo.Database, q string) (s bool, res string) {
	if len(q) != 64 {
		return false, q
	}
	tx, err := GetTx(db, q)
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
