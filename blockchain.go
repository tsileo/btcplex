package btcplex

import (
    "errors"
    "strconv"
    "fmt"
    "github.com/jmhodges/levigo"
    "github.com/pmylund/go-cache"
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

// Return the Block at the given height
func GetBlockByHeight(db *levigo.DB, blockHeight uint, fetchTx bool) (block *Block, err error) {
    hash, err := GetBlockHashByHeight(db, blockHeight)
    if err != nil {
        return
    }
    block, _ = GetBlock(db, hash, fetchTx)
    return
}

// Return block hash in the main chain given the height
func GetBlockHashByHeight(db *levigo.DB, blockHeight uint) (hash string, err error) {
    ro := levigo.NewReadOptions()
    defer ro.Close()
    blocks, _ := GetRange(db, []byte(fmt.Sprintf("bl!height!%v!", blockHeight)), []byte(fmt.Sprintf("bl!height!%v!\xff", blockHeight)))
    for _, bl := range blocks {
        blStatus, _ := db.Get(ro, []byte(fmt.Sprintf("bl!%v!main", bl.Value)))
        blStatus2, _ := strconv.Atoi(string(blStatus[:]))
        if blStatus2 == 1 {
            hash = bl.Value
            return
        }
    }
    err = errors.New("Block not found")
    return
}

// Fetch last X blocks from block start to block start - x - 1
func GetLastXBlocks(db *levigo.DB, start uint, x uint) (blocks []*Block, err error) {
    stop := start - x
    stop--
    for i := start; i > stop + 1; i-- {
        block, blockErr := GetBlockByHeight(db, uint(i), false)
        if blockErr != nil {
            err = errors.New(fmt.Sprintf("Missing block at height", i))
            return
        }
        blocks = append(blocks, block)
    }
    return
}
