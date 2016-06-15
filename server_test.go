package raftserver

import (
	"log"
	"testing"
	"time"
)

// func TestRPCServerAlive(t *testing.T) {
// 	log.Println("TestRPCServerAlive")

// 	srv := StartServer("127.0.0.1:1234", 1)
// 	if srv == nil {
// 		t.Error("Srv init error")
// 	}

// 	time.Sleep(time.Second * 5)
// }

func TestBasicElection(t *testing.T) {
	log.Println("TesMultipleRPCServerAlive")
	var SrvList []*KVRaft
	srv1 := StartClusterServers("127.0.0.1:1230", 1, []string{"127.0.0.1:1231", "127.0.0.1:1232"})
	SrvList = append(SrvList, srv1)
	srv2 := StartClusterServers("127.0.0.1:1231", 2, []string{"127.0.0.1:1230", "127.0.0.1:1232"})
	SrvList = append(SrvList, srv2)
	srv3 := StartClusterServers("127.0.0.1:1232", 3, []string{"127.0.0.1:1231", "127.0.0.1:1230"})
	SrvList = append(SrvList, srv1)
	if srv1 == nil || srv2 == nil || srv3 == nil {
		t.Error("Srv init error")
	}

	time.Sleep(time.Second * 3)

	//Check if only one leader
	leaderCount := 0
	candidate := 0
	followerCount := 0
	for _, srv := range SrvList {
		switch srv.role {
		case Leader:
			leaderCount++
		case Candidate:
			candidate++
		case Follower:
			followerCount++
		}
	}

	if candidate > 0 || followerCount != 2 || leaderCount != 1 {
		t.Error("Basic election failed:", followerCount, candidate, leaderCount)
	}
}

//func TestSingleServer(t *testing.T) {
//log.Println("TEST >>>>>TestSingleServer<<<<")

//srv := StartServer("127.0.0.1:1234", 1)

//argP := PutArgs{Key: "test1", Value: "v1"}
//var reP PutReply
//err := srv.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//argP = PutArgs{Key: "test1", Value: "v2"}
//err = srv.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//if reP.PreviousValue != "v1" {
//t.Error("Error on last value, expect v1, got :", reP.PreviousValue)
//}

//argP = PutArgs{Key: "test1", Value: "v3"}
//err = srv.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//if reP.PreviousValue != "v2" {
//t.Error("Error on last value, expect v1, got :", reP.PreviousValue)
//}

//argG := GetArgs{Key: "test1"}
//var reG GetReply
//err = srv.Get(&argG, &reG)
//if err != nil || reG.Value != "v3" {
//t.Error("Error happen on ", err, reG)
//}

//log.Println(">>", reG, err)
//log.Println("Stop server..")
//srv.kill()
//os.RemoveAll("raftexample-1")
//}

//func TestTwoServers(t *testing.T) {
//log.Println("TEST >>>>>TestTwoServers<<<<")

//var raftMsgSrvList []string
//raftMsgSrvList = append(raftMsgSrvList, "http://127.0.0.1:12379")
//raftMsgSrvList = append(raftMsgSrvList, "http://127.0.0.1:22379")

//srv1 := StartClusterServers("127.0.0.1:1234", 1, raftMsgSrvList)
//srv2 := StartClusterServers("127.0.0.1:1235", 2, raftMsgSrvList)

//argP := PutArgs{Key: "test1", Value: "v1"}
//var reP PutReply
//err := srv1.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//argP = PutArgs{Key: "test1", Value: "v2"}
//err = srv2.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//argP = PutArgs{Key: "test1", Value: "v3"}
//err = srv1.Put(&argP, &reP)
//if err != nil {
//t.Error("Error happen on ", err)
//}
//log.Println(">>", reP, err)

//argG := GetArgs{Key: "test1"}
//var reG GetReply
//err = srv1.Get(&argG, &reG)
//if err != nil || reG.Value != "v3" {
//t.Error("Error happen on ", err, reG)
//}

//err = srv2.Get(&argG, &reG)
//log.Println(">>", reG, err)
//log.Println("Stop server..")

////Kill leader test
//srv1.kill()
//os.RemoveAll("raftexample-1")
//err = srv2.Get(&argG, &reG)
//if err != nil || reG.Value != "v3" {
//t.Error("Error on kill leader happen on ", err, reG)
//}
//srv2.kill()
//os.RemoveAll("raftexample-2")
//}
