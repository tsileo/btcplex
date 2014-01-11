package btcplex

import (
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/pmylund/go-cache"
	"strconv"
)

// Return the last height processed
func GetLastHeight(db *levigo.DB) (lastHeight uint, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	lastHeightRaw, _ := db.Get(ro, []byte("last-height"))
	lastHeightInt, _ := strconv.Atoi(string(lastHeightRaw[:]))
	lastHeight = uint(lastHeightInt)
	return
}

func GetLastHeightCached(cache *cache.Cache) (lastHeight uint, err error) {
	lastheightcached, found := cache.Get("last-height")
	if found {
		lastHeight = lastheightcached.(uint)
		return
	} else {
		err = errors.New("Error getting last height from cache")
		return
	}
}

func GetLastXBlocksCached(c *cache.Cache, db *levigo.DB, start uint, x uint) (blocks []*Block, err error) {
	cachekey := fmt.Sprintf("lastxblocks%v:%v", start, x)
	cached, found := c.Get(cachekey)
	if found {
		fmt.Printf("Blocks from cache")
		return *cached.(*[]*Block), nil
	} else {

		lastblocks, err := GetLastXBlocks(db, start, x)
		c.Set(cachekey, &lastblocks, 0)
		return lastblocks, err
	}
}

// Fetch last X blocks from block start to block start - x - 1
func GetLastXBlocks(db *levigo.DB, start uint, x uint) (blocks []*Block, err error) {
	stop := start - x
	stop--
	for i := start; i > stop+1; i-- {
		block, blockErr := GetBlockByHeight(db, uint(i), false)
		if blockErr != nil {
			err = errors.New(fmt.Sprintf("Missing block at height %v", i))
			return
		}
		blocks = append(blocks, block)
	}
	return
}
