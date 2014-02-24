# BTCplex

BTCplex is an open source [Bitcoin](http://bitcoin.org/) block chain browser written in [Go](http://golang.org/), it allows you to search and navigate the [block chain](https://en.bitcoin.it/wiki/Block_chain).

**This is an early release, you might expect some bugs.**

## Requirements

- A [bitcoind](https://github.com/bitcoin/bitcoin/) instance (you can [build bitcoind in Disable-wallet mode](https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md#disable-wallet-mode))
- Go >=1.2
- Redis 2.6+
- SSDB
- 150+GB disk space / 4+GB RAM

Build btcplex database takes **1 week and few days** on a small server (dual core 1.2GHz/6GB RAM) and **28 hours** on dedicated server (i5/16GB RAM).

## Installation

Assuming you have a working Go workspace (and $GOPATH already set), Redis and SSDB already installed:

    $ git clone https://github.com/tsileo/btcplex.git
    $ cd btcplex
    $ ./build.sh
    $ cp -r config.sample.json config.json
    $ vim config.json
    $ nohup ./bin/btcplex-import > import.log&
    $ ./bin/btcplex-server


## Roadmap

- Stabilize everything for 1.0 release.

Some features that are on my TODO list:

- An easy way to monitor Bitcoin address via API (maybe using Webhooks)
- Convert BTC to fiat money easily
- An official Python module to interact with the API and offer a reliable way to track address
- An official JS lib to interact with the API
- A Watch-only addresses page
- Live notification on a unconfirmed transaction page when it actually get included in a block 
- ... (don't hesitate to request features!)

## Contribution

Contribution are welcome, see [HACKING.md](HACKING.md) and [DESIGN.md](DESIGN.md) to get started.


## Feedback / Support

You can ping me @trucsdedev/contact@btcplex.com/thomas.sileo@gmail.com if you have any feedback/issue.


## Donation

BTC: 16obt7HXb3PmyDb1wZMA2X7HYPUPHp45GB


## License

Copyright (c) 2014 Thomas Sileo

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
