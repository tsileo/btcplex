#!/bin/bash
# Build btcplex and put the binaries in ./bin
OLD_GOBIN=$GOBIN
export GOBIN="`pwd`/bin"
cp -r ./pkg $GOPATH/src/btcplex
cp -r ./cmd/* $GOPATH/src/

go get btcplex btcplex-server btcplex-prod btcplex-blocknotify
go install btcplex-server btcplex-prod btcplex-blocknotify

rm $GOPATH/src/btcplex -rf
rm $GOPATH/src/btcplex-* -rf
export GOBIN=$OLD_GOBIN