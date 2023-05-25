package main

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
	s.Raft.start()
}
