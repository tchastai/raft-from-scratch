package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type State int

const (
	Follower State = iota + 1
	Candidate
	Leader
)

type Raft struct {
	mu          sync.Mutex
	me          int
	nodes       map[int]*Node
	state       State
	currentTerm int
	votedFor    int
	voteCount   int
	log         []LogEntry
	commitIndex int
	lastApplied int
	nextIndex   []int
	matchIndex  []int
	heartbeatC  chan bool
	toLeaderC   chan bool
	Storage     Storage
	Port        string
}

func NewRaft(me int, nodes map[int]*Node, storage Storage, port string) *Raft {
	return &Raft{
		me:      me,
		nodes:   nodes,
		Storage: storage,
		Port:    port,
	}
}

func (rf *Raft) RequestVote(args VoteArgs, reply *VoteReply) error {

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return nil
	}

	if rf.votedFor == -1 {
		rf.currentTerm = args.Term
		rf.votedFor = args.CandidateID
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	}

	return nil
}

func (rf *Raft) Heartbeat(args HeartbeatArgs, reply *HeartbeatReply) error {

	if args.Term < rf.currentTerm {
		reply.Success = false
		reply.Term = rf.currentTerm
		return nil
	}

	rf.heartbeatC <- true
	if len(args.Entries) == 0 {
		reply.Success = true
		reply.Term = rf.currentTerm
		return nil
	}

	if args.PrevLogIndex > rf.getLastIndex() {
		reply.Success = false
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
		return nil
	}

	rf.log = append(rf.log, args.Entries...)
	rf.commitIndex = rf.getLastIndex()
	fmt.Println("commit index =", rf.commitIndex)
	reply.Success = true
	reply.Term = rf.currentTerm
	reply.NextIndex = rf.getLastIndex() + 1

	return nil
}

func (rf *Raft) start(kvChan <-chan KeyValue) {
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartbeatC = make(chan bool)
	rf.toLeaderC = make(chan bool)

	rpc.Register(rf)
	rpc.HandleHTTP()
	go func() {

		err := http.ListenAndServe(rf.Port, nil)
		if err != nil {
			log.Fatal("listen error: ", err)
		}
		fmt.Printf("Server is listening on port %s\n", rf.Port)
	}()

	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			switch rf.state {
			case Follower:
				select {
				case <-rf.heartbeatC:
					log.Printf("follower-%d recived heartbeat\n", rf.me)
				case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
					log.Printf("follower-%d timeout\n", rf.me)
					rf.state = Candidate
				}
			case Candidate:
				fmt.Printf("Node: %d, I'm candidate\n", rf.me)
				rf.currentTerm++
				rf.votedFor = rf.me
				rf.voteCount = 1
				go rf.broadcastRequestVote()

				select {
				case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
					rf.state = Follower
				case <-rf.toLeaderC:
					fmt.Printf("Node: %d, I'm leader\n", rf.me)
					rf.state = Leader

					rf.nextIndex = make([]int, len(rf.nodes))
					rf.matchIndex = make([]int, len(rf.nodes))
					for i := range rf.nodes {
						rf.nextIndex[i] = 1
						rf.matchIndex[i] = 0
					}

					go func() {
						i := 0
						for {
							select {
							case msg := <-kvChan:
								kv := NewKeyValue(msg.Key, msg.Value)
								rf.log = append(rf.log, LogEntry{rf.currentTerm, i, *kv})
								fmt.Printf("Raft received key=%s and value=%s\n", msg.Key, msg.Value)

								// err := rf.Storage.Set(kv.Key, kv.Value)
								// if err == nil {
								// 	rf.mu.Lock()
								// 	rf.log = append(rf.log, LogEntry{rf.currentTerm, i, *kv})
								// 	i++
								// 	rf.mu.Unlock()
								// } else {
								// 	fmt.Printf("An error occured while trying to add new element to database err=%s", err)
								// }
							}
						}
					}()

					// go func() {
					// 	i := 0
					// 	for {
					// 		i++
					// 		kv := NewKeyValue("toto"+fmt.Sprint(i), "totozaure")
					// 		rf.log = append(rf.log, LogEntry{rf.currentTerm, i, *kv})
					// 		time.Sleep(3 * time.Second)
					// 	}
					// }()
				}
			case Leader:
				rf.broadcastHeartbeat()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (rf *Raft) broadcastRequestVote() {
	var args = VoteArgs{
		Term:        rf.currentTerm,
		CandidateID: rf.me,
	}

	for i := range rf.nodes {
		go func(i int) {
			var reply VoteReply
			rf.sendRequestVote(i, args, &reply)
		}(i)
	}
}

func (rf *Raft) sendRequestVote(serverID int, args VoteArgs, reply *VoteReply) {
	client, err := rpc.DialHTTP("tcp", rf.nodes[serverID].Address)
	if err != nil {
		fmt.Printf("An error occured while connecting to nodes with ID=%d with error=%s", serverID, err)
		log.Fatal("dialing: ", err)
	}

	defer client.Close()
	client.Call("Raft.RequestVote", args, reply)

	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.state = Follower
		rf.votedFor = -1
		return
	}

	if reply.VoteGranted {
		rf.voteCount++
	}

	if rf.voteCount >= len(rf.nodes)/2+1 {
		rf.toLeaderC <- true
	}
}

func (rf *Raft) broadcastHeartbeat() {
	for i, v := range rf.nodes {
		if v.Address == rf.Port {
			continue
		}

		var args HeartbeatArgs
		args.Term = rf.currentTerm
		args.LeaderID = rf.me
		args.LeaderCommit = rf.commitIndex

		prevLogIndex := rf.nextIndex[i] - 1
		if rf.getLastIndex() > prevLogIndex {
			args.PrevLogIndex = prevLogIndex
			args.PrevLogTerm = rf.log[prevLogIndex].LogTerm
			args.Entries = rf.log[prevLogIndex:]
			log.Printf("send entries: %v\n", args.Entries)
		}

		go func(i int, args HeartbeatArgs) {
			var reply HeartbeatReply
			rf.sendHeartbeat(i, args, &reply)
		}(i, args)
	}
}

func (rf *Raft) sendHeartbeat(serverID int, args HeartbeatArgs, reply *HeartbeatReply) {
	client, err := rpc.DialHTTP("tcp", rf.nodes[serverID].Address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	defer client.Close()
	client.Call("Raft.Heartbeat", args, reply)

	if reply.Success {
		if reply.NextIndex > 0 {
			rf.nextIndex[serverID] = reply.NextIndex
			rf.matchIndex[serverID] = rf.nextIndex[serverID] - 1
		}
	} else {

		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = Follower
			rf.votedFor = -1
			return
		}
	}
}

func (rf *Raft) getLastIndex() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogIndex
}

func (rf *Raft) getLastTerm() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogTerm
}
