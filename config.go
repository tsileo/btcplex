package btcplex

import (
	"io/ioutil"
	"encoding/json"
)

// Struct holding our configuration
type Config struct {
	BitcoindBlocksPath string `json:"bitcoind_blocks_path"` 
	BitcoindRpcUrl string `json:"bitcoind_rpc_url"`
	RedisHost string `json:"redis_host"`
	MongoDbHost string `json:"mongodb_host"`
	MongoDbDb string `json:"mongodb_db"`
	LevelDbPath string  `json:"leveldb_path"`
}

// Load configuration from json file
func LoadConfig(path string) (conf *Config, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	conf = new(Config)
	json.Unmarshal(file, conf)
	return
}
