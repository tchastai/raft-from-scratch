package main

import (
	"flag"
	"strings"
	"time"
)

func main() {
	serverPort := flag.String("serverport", ":3000", "Give your listen port for server")
	raftPort := flag.String("raftport", ":3001", "Give your listen port for raft")
	cluster := flag.String("cluster", "127.0.0.1:3000", "Give cluster ip comma seperated")
	id := flag.Int("id", 1, "Give the actual node id")
	flag.Parse()

	clusters := strings.Split(*cluster, ",")

	time.Sleep(time.Second * 5)

	ns := make(map[int]*Node)
	for k, v := range clusters {
		ns[k] = NewNode(v)
	}

	storage := NewStore()
	raft := NewRaft(*id, ns, storage, *raftPort)
	serverOpts := NewServerOpts(*serverPort)
	server := NewServer(*serverOpts, *raft)
	server.Start()
	select {}

}
