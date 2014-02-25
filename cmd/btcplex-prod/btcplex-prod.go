// Process new block and unconfirmed transactions (via RPC).
package main

import (
	"log"
	"os"

	"github.com/docopt/docopt.go"

	btcplex "github.com/tsileo/btcplex/pkg"
)

func main() {
	var err error
	usage := `Process new block and unconfirmed transactions.

Usage:
  btcplex-prod [--config=<path>]
  btcplex-prod -h | --help

Options:
  -h --help     	Show this screen.
  -c <path>, --config <path>	Path to config file [default: config.json].
`

	arguments, _ := docopt.Parse(usage, nil, true, "btcplex-prod", false)

	confFile := "config.json"
	if arguments["--config"] != nil {
		confFile = arguments["--config"].(string)
	}

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %v", confFile)
	}

	conf, err := btcplex.LoadConfig(confFile)
	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}
	pool, err := btcplex.GetRedis(conf)
	if err != nil {
		log.Fatalf("Can't connect to Redis: %v", err)
	}
	ssdb, err := btcplex.GetSSDB(conf)
	if err != nil {
		log.Fatalf("Can't connect to SSDB: %v", err)
	}

	log.Println("Catching up latest block before starting")
	for {
		done := btcplex.CatchUpLatestBlock(conf, pool, ssdb)
		if done {
			break
		}	
	}
	log.Println("Catch up done!")
	
	go btcplex.ProcessNewBlock(conf, pool, ssdb)

	conn := pool.Get()
	defer conn.Close()

	// Process unconfirmed transactions (power the unconfirmed txs page/API)
    running := true
    btcplex.ProcessUnconfirmedTxs(conf, pool, &running)
}
