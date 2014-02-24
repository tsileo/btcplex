#!/bin/bash
# Build btcplex and put the binaries in ./bin
rm $GOPATH/src/btcplex -rf
cp -r ./pkg $GOPATH/src/btcplex
go get btcplex
go run cmd/btcplex-import/btcplex-import.go
