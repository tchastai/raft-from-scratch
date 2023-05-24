# Raft-from-scratch

Raft from scratch is a Raft consensus algorithm implementation for training purposes. Feel free to contribute if you want to improve it and help people to understand Raft under the hood. This implementation is not perfect for sure but enough to prove you know what you are talking about. The goal is not to create another Raft library but just to sleep with big headache tonight.

## Contributing

### Requirement

- [Go](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/)

### How to contribute

- Fork the project
- Create a new branch with a correct name as "Feature: ...."
- Push your branch and make a pull request with explanation
- Wait to be validated

## How to start

Build application

```bash
make build
```
Start applications with two nodes

```bash
./bin/raft-from-scratch --id 1 --cluster 127.0.0.1:45000,127.0.0.1:47000 --port :45000
./bin/raft-from-scratch --id 2 --cluster 127.0.0.1:45000,127.0.0.1:47000 --port :47000
```

## About Raft Algorithm

### What is Raft ?

The Raft algorithm is a consensus algorithm designed for managing replicated logs in a distributed system. It ensures that a group of nodes agree on the same log entries in a reliable and fault-tolerant manner.

In Raft, the nodes are organized into a cluster and collectively form a replicated state machine. The algorithm elects a leader among the nodes, who acts as the central point of coordination for the cluster. The leader receives client requests, appends them to its log, and replicates the log entries to other nodes in the cluster.

Raft employs a strong leader approach, where the leader handles all log replication and consistency decisions. If the leader fails or becomes disconnected, a new leader is elected through a leader election process. The leader election ensures that there is always a stable leader in the cluster.

To maintain consistency, Raft uses a log replication mechanism. When a leader receives a client request, it appends the request to its log and then sends the log entry to other nodes, asking them to replicate it. Once a majority of nodes acknowledge the replication, the leader considers the entry committed and notifies the followers. This ensures that all committed entries are consistent across the cluster.

Raft also incorporates mechanisms to handle network partitions, node failures, and rejoining nodes. It uses heartbeats to monitor node availability and detects if a node fails to respond. When a failed node recovers or a new node joins the cluster, it can catch up by requesting and receiving log entries from the leader or other nodes.

The Raft algorithm provides a straightforward and understandable approach to consensus, making it easier to implement and reason about compared to other consensus algorithms. It strikes a balance between simplicity and performance, making it suitable for a wide range of distributed systems.

### Sources
- https://eli.thegreenplace.net/2020/implementing-raft-part-2-commands-and-log-replication/ 
- https://github.com/hashicorp/raft 
- https://raft.github.io/raft.pdf



