#!/bin/bash
# Build btcplex and put the binaries in ./bin
OLD_GOBIN=$GOBIN
export GOBIN="`pwd`/bin"
go install -a ./btcplex-server ./cmd/...
export GOBIN=$OLD_GOBIN
export OLD_GOBIN=