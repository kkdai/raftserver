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
	"net/rpc"
)

type Clerk struct {
	server string
}

func MakeClerk(server string) *Clerk {
	ck := new(Clerk)
	ck.server = server
	return ck
}

func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		log.Println("[Client] Dial err:", errx)
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		log.Println("[Client] Call err:", err)
		return true
	}

	fmt.Println(err)
	return false
}

//Get
// fetch the current value for a key.
func (ck *Clerk) Get(key string) string {
	// arg := GetArgs{Key: key}
	// var reply GetReply
	// err := call(ck.server, "KVRaft.Get", &arg, &reply)
	// if err {
	// 	log.Println(reply.Err)
	// }

	// return reply.Value
	return ""
}

//Put
func (ck *Clerk) Put(key string, value string) {
	// arg := PutArgs{Key: key, Value: value}
	// var reply PutReply

	// err := call(ck.server, "KVRaft.Put", &arg, &reply)
	// if err {
	// 	log.Println(reply.Err)
	// }
}
