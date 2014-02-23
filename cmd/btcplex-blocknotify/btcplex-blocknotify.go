// Command executed by bitcoind when a new block is found,
// publish the hash over a Redis PubSub channel.
package main

import (
	"github.com/docopt/docopt.go"
	"log"
	"os"

	btcplex "github.com/tsileo/btcplex/pkg"
)

func main() {
	usage := `Callback executed when bitcoind best block changes.

Usage:
  newblock [--config=<path>] <hash>
  newblock -h | --help

Options:
  -h --help     	Show this screen.
  -c <path>, --config <path>	Path to config file [default: config.json].
`

	arguments, _ := docopt.Parse(usage, nil, true, "newblock", false)

	confFile := "config.json"
	if arguments["--config"] != nil {
		confFile = arguments["--config"].(string)
	}

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %v", confFile)
	}

	conf, _ := btcplex.LoadConfig(confFile)
	pool, _ := btcplex.GetRedis(conf)

	conn := pool.Get()
	defer conn.Close()

	conn.Do("PUBLISH", "btcplex:blocknotify", arguments["<hash>"].(string))
}
