package btcplex

import (
	"fmt"
	"labix.org/v2/mgo"
	"os"
)

func GetMongoDB() (db *mgo.Database, sess *mgo.Session, err error) {
	uri := "localhost"
	sess, err = mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	sess.SetSafe(nil)
	db = sess.DB("btcplex")

	err = db.C("blocks").EnsureIndex(mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	})
	if err != nil {
		return
	}

	err = db.C("blocks").EnsureIndex(mgo.Index{
		Key: []string{"height"},
	})
	if err != nil {
		return
	}

	err = db.C("txos").EnsureIndex(mgo.Index{
		Key: []string{"txhash", "index"},
	})
	if err != nil {
		return
	}

	err = db.C("txis").EnsureIndex(mgo.Index{
		Key: []string{"txhash", "index"},
	})
	if err != nil {
		return
	}

	err = db.C("txs").EnsureIndex(mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	})
	if err != nil {
		return
	}

	err = db.C("txs").EnsureIndex(mgo.Index{
		Key: []string{"blockhash"},
	})
	if err != nil {
		return
	}

	err = db.C("txs").EnsureIndex(mgo.Index{
		Key: []string{"blockhash", "index"},
	})
	if err != nil {
		return
	}

	err = db.C("txos").EnsureIndex(mgo.Index{
		Key: []string{"addr"},
	})
	if err != nil {
		return
	}

	err = db.C("txis").EnsureIndex(mgo.Index{
		Key: []string{"prevout.address"},
	})
	if err != nil {
		return
	}

	return
}
