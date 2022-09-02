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

// Package queue provides a queue to processes elements in FIFO order using
// a specified handler, while allowing non-blocking writes to the queue.
// The queue may be bounded in size or unbounded.
package queue

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	log "github.com/golang/glog"
)

// A Queue is a set of elements that are processed in FIFO order.  The writes
// to the queue are non-blocking and the client can read from a channel it can
// access by Receive().  The size of the queue is either bounded or unbounded.
// Note that counters such as max/count used for synchronization are managed
// under the queue lock, while counters for book keeping are handled via the
// atomic package.
type Queue struct {
	name     string           // name is a string describing the queue.
	ready    *sync.Cond       // signal indicating that the queue has elements.
	in       chan interface{} // input channel.
	out      chan interface{} // output channel, the client can access this channel via Receive().
	max      int              // maximum number of elements in a bounded queue.
	count    int              // current number of elements in a bounded queue.
	elements []interface{}    // elements in the queue.
	closed   bool             // true if the queue is closed.

	enqueueCount int64 // total number of elements enqueued queue.
	dequeueCount int64 // total number of elements processed.
}

// unbounded is the maximum size of an unbounded queue.
const unbounded = -1

// newQueue creates a new queue.
func newQueue(name string, max int) (*Queue, error) {
	return &Queue{
		name:   name,
		ready:  sync.NewCond(&sync.Mutex{}),
		max:    max,
		closed: false,
		in:     make(chan interface{}),
		out:    make(chan interface{}),
	}, nil
}

// NewUnbounded creates a new unbounded queue.
func NewUnbounded(name string) (*Queue, error) {
	return newQueue(name, unbounded)
}

// NewBounded creates a new bounded queue.
func NewBounded(name string, max int) (*Queue, error) {
	return newQueue(name, max)
}

// String returns a string representation of the queue and its state.
func (q *Queue) String() string {
	return fmt.Sprintf("%v;<Closed=%v, Max=%v, Count=%v, Enqueued=%v, Dequeued=%v>;", q.name, q.closed, q.max, q.count, q.EnqueueCount(), q.DequeueCount())
}

// EnqueueCount returns the number of elements enqueued in the queue.
func (q *Queue) EnqueueCount() int64 {
	return atomic.LoadInt64(&q.enqueueCount)
}

// DequeueCount returns the number of elements dequeued from the queue.
func (q *Queue) DequeueCount() int64 {
	return atomic.LoadInt64(&q.dequeueCount)
}

// Write writes an element to the queue if the queue is not closed and if it has space.
func (q *Queue) Write(v interface{}) error {
	err := func() error {
		q.ready.L.Lock()
		defer q.ready.L.Unlock()
		if q.closed {
			return errors.New("queue: queue closed")
		}
		if q.max != unbounded {
			if q.count >= q.max {
				return errors.New("queue: max capacity reached")
			}
			q.count++
		}
		return nil
	}()
	if err != nil {
		return err
	}
	q.in <- v
	return nil
}

// Receive returns the channel that the queue will output its elements to the client.
func (q *Queue) Receive() chan interface{} {
	return q.out
}

// Close closes the queue.
func (q *Queue) Close() {
	q.ready.L.Lock()
	defer q.ready.L.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	close(q.in)
	// No more packets are processed.
	q.elements = nil
	q.ready.Signal()
}

// Run starts processing the Queue.
func (q *Queue) Run() {
	var started sync.WaitGroup
	started.Add(2)
	// Process the in channel.
	go func() {
		log.Infof("Queue: running Queue %v", q)
		started.Done()
		for {
			v, ok := <-q.in
			if !ok {
				log.Infof("Queue: in channel closed %v", q)
				return
			}
			closed := func() bool {
				q.ready.L.Lock()
				defer q.ready.L.Unlock()
				if q.closed {
					return true
				}
				q.elements = append(q.elements, v)
				atomic.AddInt64(&q.enqueueCount, 1)
				q.ready.Signal()
				return false
			}()
			if closed {
				return
			}
		}
	}()
	// Process the out channel.
	go func() {
		started.Done()
		for {
			v, closed := func() (interface{}, bool) {
				q.ready.L.Lock()
				defer q.ready.L.Unlock()
				if !q.closed && len(q.elements) == 0 {
					q.ready.Wait()
				}
				if q.closed {
					log.Infof("Queue: Closed, closing out channel %v", q)
					close(q.out)
					return nil, true
				}
				v := q.elements[0]
				q.elements = q.elements[1:]
				return v, false
			}()
			if closed {
				return
			}
			q.out <- v
			atomic.AddInt64(&q.dequeueCount, 1)
			if q.max != unbounded {
				q.ready.L.Lock()
				q.count--
				q.ready.L.Unlock()
			}
		}
	}()
	started.Wait()
}
