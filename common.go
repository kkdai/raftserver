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

const (
	OK           = "OK"
	ErrNoKey     = "ErrNoKey"
	InvalidParam = "Invalid Parameter"
)

//AEParam : For AppendEntries
type AEParam struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []string
	LeaderCommit int
}

//AEReply :
type AEReply struct {
	Term    int
	Success bool
}

//RVParam
type RVParam struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//RVReply
type RVReply struct {
	Term        int
	VoteGranted bool
}

type Err string

type PutArgs struct {
	Key   string
	Value string
}

type PutReply struct {
	Err           Err
	PreviousValue string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Err   Err
	Value string
}
