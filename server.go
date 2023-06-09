package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	Raft   *Raft
	Http   *fiber.App
	KvChan chan KeyValue
}

func NewServer(serverOpts ServerOpts, raft Raft) *Server {
	return &Server{
		ServerOpts: &serverOpts,
		Raft:       &raft,
		KvChan:     make(chan KeyValue),
	}
}

func (s *Server) Start() {

	go s.Raft.start(s.KvChan)

	s.ListenAndServe(s.KvChan)

}

func (s *Server) ListenAndServe(kvChan <-chan KeyValue) {
	s.Http = fiber.New()
	s.Http.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})
	s.Http.Post("/set/:key-:value", func(c *fiber.Ctx) error {
		k := c.Params("key")
		v := c.Params("value")
		kv := NewKeyValue(k, v)
		s.KvChan <- *kv
		fmt.Printf("SET routes with Key=%s and Val=%s\n", k, v)
		return c.SendString("SET Routes")
	})

	s.Http.Get("/get/:key", func(c *fiber.Ctx) error {
		k := c.Params("key")
		fmt.Printf("GET routes with Key=%s\n", k)
		return c.SendString("GET Routes")
	})
	s.Http.Listen(s.Port)
}
