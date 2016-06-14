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
// limitations under the License

package raftserver

// func TestClientWithSingleServer(t *testing.T) {
// 	srv := StartServer("127.0.0.1:1234", 1)

// 	client := MakeClerk("127.0.0.1:1234")
// 	client.Put("t1", "v1")
// 	ret := client.Get("t1")
// 	log.Println("got:", ret)
// 	if ret != "v1" {
// 		t.Error("Client get error:", ret)
// 	}
// 	srv.kill()
// 	os.RemoveAll("raftexample-1")
// }
