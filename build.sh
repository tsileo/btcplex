#!/bin/bash
# Build btcplex and put the binaries in ./bin
OLD_GOBIN=$GOBIN
export GOBIN="`pwd`/bin"
cp -r ./pkg $GOPATH/src/btcplex
cp -r ./btcplex-server $GOPATH/src/

go get btcplex
go get btcplex-server
go install btcplex-server

rm ./tmp -rf
rm $GOPATH/src/btcplex -rf
rm $GOPATH/src/btcplex-server -rf
export GOBIN=$OLD_GOBIN
export OLD_GOBIN=