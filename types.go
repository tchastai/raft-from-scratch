package main

type LogEntry struct {
	LogTerm  int
	LogIndex int
	LogCMD   interface{}
}

type VoteArgs struct {
	Term        int
	CandidateID int
}

type VoteReply struct {
	Term        int
	VoteGranted bool
}

type HeartbeatArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type HeartbeatReply struct {
	Success   bool
	Term      int
	NextIndex int
}

type Node struct {
	Connect bool
	Address string
}

func NewNode(address string) *Node {
	return &Node{
		Address: address,
	}
}
