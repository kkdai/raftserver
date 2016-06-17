RaftServer: A RPC Server implement base on Raft Paper in Golang
==============

[![GoDoc](https://godoc.org/github.com/kkdai/raftserver?status.svg)](https://godoc.org/github.com/kkdai/raftserver)  [![Build Status](https://travis-ci.org/kkdai/raftserver.svg?branch=master)](https://travis-ci.org/kkdai/raftserver)
[![](https://goreportcard.com/badge/github.com/kkdai/raftserver)](https://goreportcard.com/badge/github.com/kkdai/raftserver)


What is Raft
=============

Raft is a consensus algorithm that is designed to be easy to understand. It's equivalent to Paxos in fault-tolerance and performance. The difference is that it's decomposed into relatively independent subproblems, and it cleanly addresses all major pieces needed for practical systems. We hope Raft will make consensus available to a wider audience, and that this wider audience will be able to develop a variety of higher quality consensus-based systems than are available today. (quote from [here](https://raft.github.io/))



Installation and Usage
=============

Install
---------------
``` 
go get github.com/kkdai/raftserver
```

Usage
---------------
```
	//Start three raft servers
	srv1 := StartClusterServers("127.0.0.1:1230", 1, []string{"127.0.0.1:1231", "127.0.0.1:1232"})
	srv2 := StartClusterServers("127.0.0.1:1231", 2, []string{"127.0.0.1:1230", "127.0.0.1:1232"})
	srv3 := StartClusterServers("127.0.0.1:1232", 3, []string{"127.0.0.1:1231", "127.0.0.1:1230"})

	//Need sleep a while for leader election.
	time.Sleep(time.Second * 2)

	//Ask who is leader
	fmt.Println(srv1.WhoAreYou())
	//Leader (It depends on election result)		
	
```



Inspired By
---------------

- [In Search of an Understandable Consensus Algorithm](https://ramcloud.stanford.edu/raft.pdf)


Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

etcd is under the Apache 2.0 [license](LICENSE). See the LICENSE file for details.