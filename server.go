// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raftserver

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"syscall"
)

const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

type KVRaft struct {
	mu         sync.Mutex
	l          net.Listener
	me         int
	dead       bool // for testing
	unreliable bool // for testing

	// Raft Persistent State
	currentTerm int
	voteFor     int
	log         *LogData
	role        Role

	// Raft Volatile state (All Server)
	commitIndex int
	lastApplied int

	// Raft Volatile state (Leader)
	nextIndex  int
	matchIndex int
}

func NewKVRaft() *KVRaft {
	k := new(KVRaft)
	k.log = NewLogData()
	k.role = Follower
	return k
}

//change role
func (kv *KVRaft) changeRole(r Role) {
	log.Println("[ROLE] Server:", kv.me, " change from ", kv.role, " to ", r)
	kv.role = r
}

//AppendEntries :RPC call to server to update Log Entry
func (kv *KVRaft) AppendEntries(args *AEParam, reply *AEReply) error {
	DPrintf("[AppendEntries] args=%v", args)
	//Reply false if term small than current one
	if args.Term < kv.currentTerm {
		reply.Success = false
		//apply new term and change back to follower
		kv.currentTerm = args.Term
		kv.changeRole(Follower)
		return errors.New("term is smaller than current")
	}

	//Reply false if log doesn;t contain an entry at previous
	if exist, index := kv.log.ContainIndex(args.PrevLogTerm, args.PrevLogIndex); !exist {
		reply.Success = false
		return errors.New("No contain previous index or term term =>" + fmt.Sprint(args))
	} else {
		//If eixst a conflict value, trim it to the end
		if data, err := kv.log.Get(index); err != nil {
			log.Println("Get data failed on exist index, racing")
		} else {
			if data.Term != args.Term || data.Data != args.Entries[0] {
				kv.log.TrimRight(index)
			}

		}
	}

	//Append log to storage
	var inLogs []Log
	for _, v := range args.Entries {
		inLogs = append(inLogs, Log{Term: args.Term, Data: v})
	}

	//Appen data
	kv.log.Append(inLogs)

	//Update reply
	reply.Term = kv.currentTerm
	reply.Success = true

	//update commitIndex
	if args.LeaderCommit > kv.commitIndex {
		if args.LeaderCommit <= (kv.log.Length() - 1) {
			kv.commitIndex = args.LeaderCommit
		} else {
			kv.commitIndex = (kv.log.Length() - 1)
		}
	}

	// update other
	kv.currentTerm = args.Term
	kv.lastApplied = kv.log.Length() - 1 //update latest state machine index
	if kv.lastApplied < kv.commitIndex {
		kv.lastApplied = kv.commitIndex
	}

	//TODO  Heart beat
	return nil
}

func (kv *KVRaft) RequestVote(args *RVParam, reply *RVReply) error {
	DPrintf("[RequestVote] args=%v", args)
	//Reply false if term small than current one
	if args.Term < kv.currentTerm {
		reply.VoteGranted = false
		kv.currentTerm = args.Term
		kv.changeRole(Follower)
		return errors.New("term is smaller than current")
	}

	if kv.voteFor == 0 && args.LastLogIndex >= kv.lastApplied && args.LastLogTerm >= kv.currentTerm {
		kv.voteFor = args.CandidateId
		reply.Term = kv.currentTerm
		reply.VoteGranted = true
		return nil
	}

	return errors.New("Vote failed:" + fmt.Sprint("%v -> %v", kv.lastApplied, *args))
}

func (k *KVRaft) serverLoop() {

}

// tell the server to shut itself down.
// please do not change this function.
func (kv *KVRaft) kill() {
	DPrintf("Kill(%d): die\n", kv.me)
}

func StartServer(rpcPort string, me int) *KVRaft {
	return startServer(rpcPort, me, []string{rpcPort}, false)
}

func StartClusterServers(rpcPort string, me int, cluster []string) *KVRaft {
	return startServer(rpcPort, me, cluster, false)
}

func StarServerJoinCluster(rpcPort string, me int) *KVRaft {
	return startServer(rpcPort, me, []string{rpcPort}, true)
}

func startServer(serversPort string, me int, cluster []string, join bool) *KVRaft {
	kv := NewKVRaft()
	rpcs := rpc.NewServer()
	rpcs.Register(kv)
	go kv.serverLoop()

	DPrintf("[server] %d port:%s", me, serversPort)
	l, e := net.Listen("tcp", serversPort)
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	kv.l = l

	go func() {
		for kv.dead == false {
			conn, err := kv.l.Accept()
			if err == nil && kv.dead == false {
				if kv.unreliable && (rand.Int63()%1000) < 100 {
					// discard the request.
					conn.Close()
				} else if kv.unreliable && (rand.Int63()%1000) < 200 {
					// process the request but force discard of reply.
					c1 := conn.(*net.UnixConn)
					f, _ := c1.File()
					err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
					if err != nil {
						fmt.Printf("shutdown: %v\n", err)
					}
					go rpcs.ServeConn(conn)
				} else {
					go rpcs.ServeConn(conn)
				}
			} else if err == nil {
				conn.Close()
			}
			if err != nil && kv.dead == false {
				fmt.Printf("KVRaft(%v) accept: %v\n", me, err.Error())
				kv.kill()
			}
		}
	}()

	return kv
}
