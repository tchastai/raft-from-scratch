package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

type ServerOpts struct {
	Port string
}

func NewServerOpts(port string) *ServerOpts {
	return &ServerOpts{
		Port: port,
	}
}

type Server struct {
	*ServerOpts

	Raft *Raft
}

func NewServer(serverOpts ServerOpts, raft Raft) *Server {
	return &Server{
		ServerOpts: &serverOpts,
		Raft:       &raft,
	}
}

func (s *Server) Start() {
	rpc.Register(s.Raft)
	rpc.HandleHTTP()
	go func() {
		err := http.ListenAndServe(s.Port, nil)
		if err != nil {
			log.Fatal("listen error: ", err)
		}
		fmt.Printf("Server is listening on port %s\n", s.Port)
	}()
	s.Raft.start()
}
