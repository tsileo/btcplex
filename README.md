# BTCplex

BTCplex is an open source [Bitcoin](http://bitcoin.org/) block chain browser written in [Go](http://golang.org/), it allows you to search and navigate the [block chain](https://en.bitcoin.it/wiki/Block_chain).


## Roadmap

Some features that are on my TODO list:

- An easy way to monitor Bitcoin address via API (maybe using Webhooks)
- Convert BTC to USD easily
- An official Python and JS API wrapper
- A Watch-only addresses page
- Live notification on a unconfirmed transaction page when it actually get included in a block 

## Architecture

[Go](http://golang.org/) for the server-side ([Martini](http://martini.codegangsta.io/) for the webapp), [Redis](http://redis.io/) for temporary data, [SSDB](https://github.com/ideawu/ssdb) (Backed by [LevelDB](https://code.google.com/p/leveldb/)) for persistent data.

The initial import is still slow, I've spent already spent a lot of time trying to improve it (1+ week on a small dual-core/6GB RAM server at block 247000), and the database takes around 69GB (at block 247000). I think the minimal requirement are 100+GB HD and 4+GB RAM.

### Unconfirmed transactions

Bitcoind memory pool is synced every 500ms in Redis, along with every transaction.

The most recent memory pool sync is stored in a sorted set (with time as score), allowing them to be "replayed" via SSE on the unconfirmed transactions page.
Each unconfirmed transaction is stored as JSON in a key ``btcplex:utx:%v`` (hash), the key is destroyed when it get removed from the memory pool.
The sync is performed by keeping two sets: ``btcplex:rawmempool:%v`` (unix time):

- one containing the previous state of memory pool (500ms ago)
- the current state of memory pool

The diff of the two sets is computed, and old unconfirmed transactions are removed. 


### Available keys

I will try to keep an updated list of how data is stored in SSDB by data type.


#### Strings

Blocks, transactions, TxIns, TxOuts are stored in JSON format (SSDB support Redis protocol but it doesn't support MULTI, so I can't use hashes if I can't retrieve multiple hashes in one call, also, some JSON objects are nested, so I stick to JSON.

- ``block:height:%v`` (height) -> Contains the hash for the given height
- ``block:%v`` (hash) -> Block data in JSON format
- ``tx:%v`` (hash) -> Transaction data in JSON format
- ``txi:%v:%v`` (hash, index) -> TxIn data in JSON format
- ``txo:%v:%v`` (hash, index) -> TxOut data in JSON format
- ``txo:%v:%v:spent`` (hash, index) -> Spent data in JSON format


#### Hashes

BTCplex keeps one hashes per address (``addr:%v:h`` (address)) containing TotalSent, TotalReceived. 

Hash keys for address data:

- ``ts`` -> TotalSent
- ``tr`` -> TotalReceived


#### Sorted Sets

BTCplex keeps three sorted set per address (``addr:%v`` (address), ``addr:%v:received`` (address), ``addr:%v:sent`` (address)) containing Tx hash for every transaction involving the address sorted by BlockTime.

It also store one sorted for each block containing transaction references sorted by index (``block:%v:txs`` (hash)).


### Webapp

- [Martini](http://martini.codegangsta.io/)


### Backend

I tried a lot of databases ([LevelDB](https://code.google.com/p/leveldb/), [RethinkDB](http://rethinkdb.com/), [MongoDB](http://mongodb.org/), and [Ardb](https://github.com/yinqiwen/ardb)), and [SSDB](https://github.com/ideawu/ssdb) was the faster, I didn't wanted to use Redis because it would need a lot of RAM and I wanted persistent storage but since SSDB is a drop-in replacement for Redis, you can only use Redis if you prefer. If you choose to use SSDB, [Redis](http://redis.io/) is still needed for PubSub (used for SSE). [LevelDB](https://code.google.com/p/leveldb/) is used for caching during the initial import.

- [Redis](http://redis.io/)
- [SSDB](https://github.com/ideawu/ssdb) or Redis
- [LevelDB](https://code.google.com/p/leveldb/)


## Donation

BTC: 16obt7HXb3PmyDb1wZMA2X7HYPUPHp45GB


## Feedback / Support

You can ping me @trucsdedev/contact@btcplex.com/thomas.sileo@gmail.com if you have any feedback/issue.


## License

Copyright (c) 2014 Thomas Sileo

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.