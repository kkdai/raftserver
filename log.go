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
	"log"
)

type Log struct {
	Term int
	Data string
}

type LogData struct {
	data []Log
}

func NewLogData() *LogData {
	return new(LogData)
}

//ContainIndex :Reply true if term and index exist in current log, return index if exist
func (l *LogData) ContainIndex(term int, index int) (bool, int) {
	for k, v := range l.data {
		if v.Term == term && k == index {
			return true, k
		}
	}
	return false, -1
}

//Get :Get data by index
//TrimRight :Trim data after specfic index
func (l *LogData) TrimRight(index int) {
	if len(l.data) < index {
		return
	}

	l.data = l.data[index:]
}

//Get :Get data by index
func (l *LogData) Get(index int) (Log, error) {
	if l.Length() < index {
		return Log{}, errors.New("No data")
	}

	return l.data[index], nil
}

//Append :Append data to last if not exist
func (l *LogData) Append(in []Log) error {
	var lastLog Log

	//empty for heart beat, just return nil
	if len(in) == 0 {
		log.Println("Empty entry for heart beat")
		return nil
	}

	if len(l.data) > 0 {
		lastLog = l.data[len(l.data)-1]
	}

	if lastLog.Term == in[0].Term && lastLog.Data == in[0].Data {
		return errors.New("Append failed, exist")
	}

	for _, v := range in {
		l.data = append(l.data, v)
	}
	return nil
}

//Length :Return the data length
func (l *LogData) Length() int {
	return len(l.data)
}
