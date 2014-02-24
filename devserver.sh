#!/bin/bash
# Build btcplex and put the binaries in ./bin
rm $GOPATH/src/btcplex -rf
cp -r ./pkg $GOPATH/src/btcplex
cp -r ./btcplex-server $GOPATH/src/btcplex-server
go get btcplex btcplex-server
go run btcplex-server/btcplex-server.go
