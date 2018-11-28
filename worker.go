package workerpool

import (
	"fmt"
)

// Worker provides a standard worker interface for the worker pool
type Worker interface {
	start(pool WorkerPool)
	stop(pool WorkerPool)
}

// NaiveWorker is a naive implementation of the Worker interface
type NaiveWorker struct {
	id   uint64
	temp bool
}

func (w *NaiveWorker) start(pool WorkerPool) {
	go func() {
		pool.registerWorker()
		for work := range *pool.workQueue() {
			work()
			if w.temp {
				break
			}
		}
		pool.unregisterWorker()
		fmt.Printf("worker %d stopping\n", w.id)
	}()
}
