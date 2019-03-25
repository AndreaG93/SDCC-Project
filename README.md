# SDCC-Project-WorkerNode
This repository is a part of an university project. It contains a Go implementation of a "WorkerNode".

## Getting Started

#### Install library dependencies

1. `go get -u github.com/aws/aws-sdk-go/...`

2. `go get -u github.com/samuel/go-zookeeper/zk`

3. `go get -u github.com/lni/dragonboat`

cd $GOPATH/src/github.com/lni/dragonboat

make install-rocksdb-ull

go get -u github.com/hashicorp/raft
github.com/hashicorp/raft-mdb