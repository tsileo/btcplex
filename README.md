# BTCplex

BTCplex is an open source [Bitcoin](http://bitcoin.org/) block chain browser written in [Go](http://golang.org/), it allows you to search and navigate the [block chain](https://en.bitcoin.it/wiki/Block_chain).

## Architecture

## Available keys

I will try to keep an updated list of how data is stored in SSDB by data type.

### Strings

Blocks, transactions, TxIns, TxOuts are stored in JSON format.

- ``block:height:%v`` (height) -> Contains the hash for the given height
- ``block:%v`` (hash) -> Block data in JSON format
- ``tx:%v`` (hash) -> Transaction data in JSON format
- ``txi:%v:%v`` (hash, index) -> TxIn data in JSON format
- ``txo:%v:%v`` (hash, index) -> TxOut data in JSON format
- ``txo:%v:%v:spent`` (hash, index) -> Spent data in JSON format


### Hashes

BTCplex keeps one hashes per address (``addr:%v:h`` (address)) containing ReceivedCnt, SentCnt, TotalSent, TotalReceived. 

Hash keys for address data:

- ``rc`` -> ReceivedCnt
- ``sc`` -> SentCnt
- ``ts`` -> TotalSent
- ``tr`` -> TotalReceived

### Sorted Sets

BTCplex keeps one sorted set per address (``addr:%v`` (address)) containing keys of every TxIn/TxOut sorted by BlockTime.

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

## License

Copyright (c) 2014 Thomas Sileo

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
