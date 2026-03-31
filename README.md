# FluxKV

A lightweight distributed key-value store built in Go.

## Overview

FluxKV is a distributed key-value store that demonstrates core distributed systems concepts such as:

* Consistent hashing (data partitioning)
* Peer-to-peer request routing
* Decentralized replication
* HTTP-based inter-node communication

Unlike traditional systems, FluxKV does **not rely on a leader node** — each node independently handles requests and replication.

## Architecture

* **API Layer** – Handles client requests (PUT, GET, DELETE)
* **Cluster Layer** – Routes requests using consistent hashing
* **Replication Layer** – Nodes replicate data to peers
* **Storage Layer** – In-memory key-value store

## Features

* Distributed key-value storage
* Consistent hashing for partitioning
* Request forwarding between nodes
* Decentralized (leaderless) replication
* Multi-node cluster support

## Replication Model

FluxKV uses a **peer-to-peer replication model**:

1. A client sends a request to any node
2. The node determines the responsible node using hashing
3. The responsible node processes the request
4. The node replicates the data to other nodes in the cluster

There is **no central leader** — all nodes participate equally.


## Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/FluxKV.git
cd FluxKV
```

### 2. Run a 3-node cluster

Open **3 terminals**:

```bash
PORT=8081 NODE_ID=node1 \
PEERS=localhost:8082,localhost:8083 \
go run cmd/server/main.go
```

```bash
PORT=8082 NODE_ID=node2 \
PEERS=localhost:8081,localhost:8083 \
go run cmd/server/main.go
```

```bash
PORT=8083 NODE_ID=node3 \
PEERS=localhost:8081,localhost:8082 \
go run cmd/server/main.go
```

## Usage

### PUT

```bash
curl -X POST localhost:8081/put \
  -H "Content-Type: application/json" \
  -d '{"key":"foo","value":"bar"}'
```

### GET

```bash
curl "localhost:8082/get?key=foo"
```

### DELETE

```bash
curl -X DELETE "localhost:8083/delete?key=foo"
```

## How It Works

* Any node can receive a request
* Consistent hashing determines which node owns the key
* Requests are forwarded if needed
* The responsible node stores the data
* Data is replicated to peer nodes

## Limitations

* Eventual consistency (no strong guarantees)
* No quorum-based writes
* No failure detection
* No conflict resolution
* In-memory storage only

## Future Improvements

* Quorum-based replication (like Dynamo/Cassandra)
* Read repair and anti-entropy
* Persistent storage (WAL)
* Failure detection (gossip)
* Conflict resolution (vector clocks / CRDTs)

## Learning Goals

FluxKV explores:

* Distributed system design
* Data partitioning with consistent hashing
* Decentralized replication
* Handling partial failures
