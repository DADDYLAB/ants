// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ants

import (
	"time"
	"fmt"
	"errors"
	"runtime"
)

// Worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type Worker struct {
	// Worker now running job Tag
	Tag string

	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan *Job

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *Worker) run() {
	go func() {
		// Panic
		defer w.recover()

		for f := range w.task {
			w.Tag = f.Tag
			if f == nil {
				w.pool.decRunning()
				return
			}
			f.F()
			w.pool.putWorker(w)
		}
	}()
}

func (w *Worker) recover() {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		if len(w.pool.panicHandlers) > 0 {
			for _, h := range w.pool.panicHandlers {
				h.Handle(errors.New(fmt.Sprintf("%v ，Stack: %s", err, buf)), w)
			}
		}
		w.pool.decRunning()
		w.pool.DeleteJob(w.Tag)
	}
}
