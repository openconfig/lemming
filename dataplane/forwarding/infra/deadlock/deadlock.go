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

// Package deadlock is a set of simple utilities meant to detect deadlocks
// based on timeouts.
package deadlock

import (
	"sync"
	"time"

	log "github.com/golang/glog"
)

// Timeout is the default time after which a function is assumed to have deadlock.
const Timeout = 10 * time.Minute

// A Timer is used to time a function, and declare a deadlock if it
// does not return in the specified amount of time.
type Timer struct {
	stop chan bool // channel used to stop the timer
}

// NewTimer returns a new timer that declares a deadlock if it is not stopped
// within the specified duration.
func NewTimer(timeout time.Duration, desc string) *Timer {
	ch := make(chan bool)
	go func() {
		timer := time.NewTimer(timeout)
		select {
		case <-timer.C:
			log.Errorf("TimedExecute failed for %v, timeout after %v. ", desc, timeout)

		case <-ch:
			timer.Stop()
		}
	}()
	return &Timer{
		stop: ch,
	}
}

// Stop stops a running timer.
func (t *Timer) Stop() {
	t.stop <- true
}

// A Monitor is used to monitor a goroutine and ensure it completes
// various sections of code (called monitored code) within the specified
// timeout. To use this, create a new Monitor type for each goroutine, and
// modify the goroutine to record when it enters or exists monitored code.
type Monitor struct {
	exit chan bool // channel used to indicating that the goroutine is exiting the monitored code

	runningMu   sync.Mutex
	running     bool   // true if the goroutine is executing its monitored code
	runningName string // describes the monitored code currently executing
}

// Enter records that the goroutine is entering a monitored section of code.
func (t *Monitor) Enter(desc string) {
	t.runningMu.Lock()
	t.running = true
	t.runningName = desc
	t.runningMu.Unlock()
}

// Exit records that the goroutine is exiting a monitored section of code.
func (t *Monitor) Exit() {
	t.runningMu.Lock()
	t.running = false
	t.runningName = ""
	t.runningMu.Unlock()
	t.exit <- true
}

// IsRunning returns true if the goroutine is running a monitored code.
// It also returns a description of the monitored code being executed.
func (t *Monitor) IsRunning() (bool, string) {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	return t.running, t.runningName
}

// NewMonitor returns a new monitor for the goroutine.
func NewMonitor(timeout time.Duration, name string) *Monitor {
	m := &Monitor{
		exit: make(chan bool),
	}

	timer := time.NewTimer(timeout)
	go func() {
		for {
			select {
			case <-timer.C:
				if running, desc := m.IsRunning(); running {
					log.Errorf("deadlock: Timeout detected for %v in %v, timeout %v", name, desc, timeout)
				}
			case <-m.exit:
			}
			timer.Reset(timeout)
		}
	}()
	return m
}
