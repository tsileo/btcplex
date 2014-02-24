# Design

[Go](http://golang.org/) for the server-side ([Martini](http://martini.codegangsta.io/) for the webapp), [Redis](http://redis.io/) for temporary data, [SSDB](https://github.com/ideawu/ssdb) (Backed by [LevelDB](https://code.google.com/p/leveldb/)) for persistent data.

The initial import is still slow, I already spent a lot of time trying to improve it (it took 1+ week on a small dual-core/6GB RAM server at block 247000), but less than 2 days on a more decent server (i5/16GB RAM).

## Unconfirmed transactions

Bitcoind memory pool is synced every 1s in Redis, along with every transactions.

The most recent memory pool sync is stored in a sorted set (with time as score, in ``btcplex:rawmempool``), allowing them to be "replayed" via SSE on the unconfirmed transactions page.
Each unconfirmed transaction is stored as JSON in a key ``btcplex:utx:%v`` (hash), the key is destroyed when it get removed from the memory pool.
The sync is performed by keeping two sets: ``btcplex:rawmempool:%v`` (unix time):

- one containing the previous state of memory pool (500ms ago)
- the current state of memory pool

The diff of the two sets is computed, and old unconfirmed transactions are removed. 

## New block

BTCplex relies on ``bitcoind`` blocknotify callback, each time the best block changes, it will be processed (via the RPC API) and immediately available. 

## Available keys in SSDB

I will try to keep an updated list of how data is stored in SSDB by data type.

### Strings

Blocks, transactions, TxIns, TxOuts are stored in JSON format (SSDB support Redis protocol but it doesn't support MULTI, so I can't use hashes if I can't retrieve multiple hashes in one call, also, some JSON objects are nested, so I stick to JSON.

- ``block:height:%v`` (height) -> Contains the hash for the given height
- ``block:%v`` (hash) -> Block data in JSON format
- ``block:%v`` (hash) -> Block data along with its transactions in JSON format
- ``tx:%v`` (hash) -> Transaction data in JSON format
- ``txi:%v:%v`` (hash, index) -> TxIn data in JSON format
- ``txo:%v:%v`` (hash, index) -> TxOut data in JSON format
- ``txo:%v:%v:spent`` (hash, index) -> Spent data in JSON format
- ``btcplex:utx:%v`` (hash) -> Unconfirmed transaction (with TxOuts/TxIns) in JSON format


### Hashes

BTCplex keeps one hash per address (``addr:%v:h`` (address)) containing TotalSent, TotalReceived. 

Hash keys for address data:

- ``ts`` -> TotalSent
- ``tr`` -> TotalReceived

And one hash per block (``block:%v:h`` (hash)) containing the following keys:

- ``main`` -> Boolean, false if the block is orphaned
- ``next`` -> Hash of the next block, if any
- ``parent`` -> Hash of the previous block


### Sorted Sets

BTCplex keeps three sorted set per address (``addr:%v`` (address), ``addr:%v:received`` (address), ``addr:%v:sent`` (address)) containing Tx hash for every transaction involving the address sorted by BlockTime.

It also store one sorted for each block containing transaction references sorted by index (``block:%v:txs`` (hash)).

Bitcoind memory pool is "synced" in a sorted set: ``btcplex:rawmempool``.


### Webapp

- [Martini](http://martini.codegangsta.io/)


### Server-Sent Events

One goroutine per server-sent events, broadcasted (using [bcast](https://github.com/grafov/bcast)) to multiple channels efficiently.

List of Redis PubSub channels:

- btcplex:utxs -> Rely unconfirmed transactions in JSON format
- btcplex:blocknotify -> Send best block hash when it changes (blocknotify callback)


## Backend notes

I tried a lot of databases ([LevelDB](https://code.google.com/p/leveldb/), [RethinkDB](http://rethinkdb.com/), [MongoDB](http://mongodb.org/), and [Ardb](https://github.com/yinqiwen/ardb)), and [SSDB](https://github.com/ideawu/ssdb) was the faster, I didn't wanted to use Redis because it would need a lot of RAM and I wanted persistent storage but since SSDB is a drop-in replacement for Redis, you can only use Redis if you prefer. If you choose to use SSDB, [Redis](http://redis.io/) is still needed for PubSub (used for SSE). [LevelDB](https://code.google.com/p/leveldb/) is used for caching during the initial import.

- [Redis](http://redis.io/)
- [SSDB](https://github.com/ideawu/ssdb) or Redis
- [LevelDB](https://code.google.com/p/leveldb/)

## Custom bitcoind

You can use a custom bitcoind, that fetch directly the previous tx and return directory for each txin, the address and value, allowing new blocks to be processed faster (making less RPC request)

// TODO(tsileo) provides a real patch

in /bitcoin/src/rpcrawtransaction.cpp:

	in.push_back(Pair("txid", txin.prevout.hash.GetHex()));
    in.push_back(Pair("vout", (boost::int64_t)txin.prevout.n));
    //Modif
    CTransaction txPrev;
    uint256 hashBlock;
    GetTransaction(txin.prevout.hash, txPrev, hashBlock);  // get the vin's previous transaction 
    CTxDestination source;
    if (ExtractDestination(txPrev.vout[txin.prevout.n].scriptPubKey, source))  // extract the destination of the previous transaction's vout[n]
    {
        CBitcoinAddress addressSource(source);              // convert this to an address
        in.push_back(Pair("address", addressSource.ToString())); // add the address to the returned object
        in.push_back(Pair("value", ValueFromAmount(txPrev.vout[txin.prevout.n].nValue))); 
    }
    //End
    Object o;
