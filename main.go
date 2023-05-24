package main

import (
	"flag"
	"strings"
	"time"
)

func main() {

	port := flag.String("port", ":3000", "Give your listen port for server")
	cluster := flag.String("cluster", "127.0.0.1:3000", "Give cluster ip comma seperated")
	id := flag.Int("id", 1, "Give the actual node id")
	flag.Parse()

	clusters := strings.Split(*cluster, ",")

	time.Sleep(time.Second * 5)

	ns := make(map[int]*node)
	for k, v := range clusters {
		ns[k] = newNode(v)
	}

	storage := NewStore()
	raft := NewRaft(*id, ns, storage)
	serverOpts := NewServerOpts(*port)
	server := NewServer(*serverOpts, *raft)
	server.Start()
	select {}

}
