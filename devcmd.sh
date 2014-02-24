#!/bin/bash
# Copy packages in $GOPATH, and run the given command directly
# e.g. ./devcmd server to run ./btcplex-server
rm $GOPATH/src/btcplex -rf
rm $GOPATH/src/btcplex-* -rf
cp -r ./pkg $GOPATH/src/btcplex
cp -r ./cmd/btcplex-$1 $GOPATH/src/btcplex-$1
go get btcplex btcplex-$1
if [ $? -eq 0 ]; then
    go run cmd/btcplex-$1/btcplex-$1.go
fi