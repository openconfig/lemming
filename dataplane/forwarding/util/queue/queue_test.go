// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// TestQueue tests operations on queue of various sizes.
func TestQueue(t *testing.T) {
	tests := []struct {
		size     int // size of the queue
		write    int // number of writes to attempt
		read     int // number of reads to attempt
		wantRead int // number of successful read expected
	}{
		{size: 10, write: 10, read: 10, wantRead: 10},
		{size: 10, write: 10, read: 5, wantRead: 5},
		{size: 5, write: 10, read: 5, wantRead: 5},
		{size: unbounded, write: 10, read: 10, wantRead: 10},
		{size: unbounded, write: 10, read: 5, wantRead: 5},
	}

	for id, test := range tests {
		c, err := newQueue(fmt.Sprintf("Test %v", id), test.size)
		if err != nil {
			t.Errorf("%d: Unexpected error in Queue creation %v.", id, err)
			continue
		}
		c.Run()

		for i := 0; i < test.write; i++ {
			wantNil := true
			if test.size != unbounded && i >= test.size {
				wantNil = false
			}
			err := c.Write(i)
			if errNil := err == nil; errNil != wantNil {
				t.Errorf("%v(%v): got %v, want nil %v", i, c, err, wantNil)
			}
		}

		readCnt := 0
		errCnt := 0
		for i := 0; i < test.read; i++ {
			if v, ok := <-c.Receive(); ok {
				readCnt++
				if v != i {
					t.Errorf("%v: got %v, want %v", c, v, i)
				}
			} else {
				errCnt++
			}
		}
		c.Close()
		if readCnt != test.wantRead {
			t.Errorf("%v: read %v, want %v", id, readCnt, test.wantRead)
		}
		if errCnt != test.read-test.wantRead {
			t.Errorf("%v: read fail %v, want %v", id, errCnt, test.read-test.wantRead)
		}
		if err := c.Write(0); err == nil {
			t.Error("Write() on closed got nil, want err")
		}
		// retry to close
		c.Close()
	}
}

// TestQueueConcurrent tests operations on queue of various sizes where many
// writers write concurrently. It waits until all writes are done before start
// reading, to avoid non determinism in the number of successful read / write
// caused by readers and writers racing when the queue size is smaller.
func TestQueueConcurrent(t *testing.T) {
	tests := []struct {
		size      int // size of the queue
		write     int // number of writes to attempt
		writeWant int // number of writes expected to succeed
		read      int // number of reads to attempt
		readWant  int // number of reads expected to succeed
	}{
		{size: 10, write: 10, writeWant: 10, read: 10, readWant: 10},
		{size: 10, write: 10, writeWant: 10, read: 5, readWant: 5},
		{size: 10, write: 5, read: 10, writeWant: 5, readWant: 5},
		{size: 5, write: 10, read: 10, writeWant: 5, readWant: 5},
		{size: 5, write: 10, read: 5, writeWant: 5, readWant: 5},
		{size: 5, write: 5, read: 10, writeWant: 5, readWant: 5},
		{size: unbounded, write: 10, read: 10, writeWant: 10, readWant: 10},
		{size: unbounded, write: 10, read: 5, writeWant: 10, readWant: 5},
		{size: unbounded, write: 5, read: 10, writeWant: 5, readWant: 5},
	}

	for id, test := range tests {
		c, err := newQueue(fmt.Sprintf("Test %v", id), test.size)
		if err != nil {
			t.Errorf("%d: Unexpected error in Queue creation %v.", id, err)
			continue
		}
		c.Run()

		var writeCnt, readCnt int32
		// key: element that is written, value: count.
		writtenElements := make(map[int]int)
		// key: element that is read, value: count.
		readElements := make(map[int]int)

		// wait for all write operations
		var waitWrite sync.WaitGroup
		var l sync.Mutex
		waitWrite.Add(test.write)
		for i := 0; i < test.write; i++ {
			go func(v int) {
				defer waitWrite.Done()
				if err := c.Write(v); err == nil {
					atomic.AddInt32(&writeCnt, 1)
					l.Lock()
					writtenElements[v]++
					l.Unlock()
				}
			}(i)
		}
		waitWrite.Wait()

		// all read should start and # of successful read should finish
		// before closing the queue.
		var start, all, toFinish sync.WaitGroup
		start.Add(test.read)
		all.Add(test.read)
		toFinish.Add(test.readWant)
		for i := 0; i < test.read; i++ {
			go func(idx int) {
				defer all.Done()
				start.Done()
				if v, ok := <-c.Receive(); ok {
					atomic.AddInt32(&readCnt, 1)
					l.Lock()
					readElements[v.(int)]++
					l.Unlock()
					toFinish.Done()
				}
			}(i)
		}
		start.Wait()
		toFinish.Wait()
		c.Close()
		all.Wait()
		if writeCnt != int32(test.writeWant) {
			t.Errorf("%v: write done: %v, want: %v", id, writeCnt, test.writeWant)
		}
		if readCnt != int32(test.readWant) {
			t.Errorf("%v: read done: %v, want: %v", id, readCnt, test.readWant)
		}
		// verify values that are read match to the values that are written.
		for k, v := range readElements {
			if c, ok := writtenElements[k]; !ok || c != v {
				t.Errorf("read unwritten item %v", k)
			} else {
				writtenElements[k] -= v
				if writtenElements[k] == 0 {
					delete(writtenElements, k)
				}
			}
		}
		sum := 0
		for _, v := range writtenElements {
			sum += v
		}
		if want := (test.writeWant - test.readWant); sum != want {
			t.Errorf("remaining item: %v, want :%v", sum, want)
		}
	}
}
