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
	log         []string

	// Raft Volatile state (All Server)
	commitIndex int
	lastApplied int

	// Raft Volatile state (Leader)
	nextIndex  int
	matchIndex int
}

func (kv *KVRaft) AppendEntries(args *AEParam, reply *AEReply) error {
	DPrintf("[AppendEntries]", args)
	return nil
}

func (kv *KVRaft) RequestVote(args *RVParam, reply *RVReply) error {
	DPrintf("[RequestVote]", args)

	return nil
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
	kv := new(KVRaft)
	rpcs := rpc.NewServer()
	rpcs.Register(kv)
	go kv.serverLoop()

	DPrintf("[server] ", me, " ==> ", serversPort)
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
