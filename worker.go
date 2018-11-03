package workerpool

import (
	"fmt"
	"sync"
)

// Worker provides a standard worker interface for the worker pool
type Worker interface {
	start(pool WorkerPool)
	stop(pool WorkerPool)
}

// NaiveWorker is a naive implementation of the Worker interface
type NaiveWorker struct {
	quitChan     chan bool
	id           uint64
	temp         bool
	shutdownOnce sync.Once
}

func (w *NaiveWorker) start(pool WorkerPool) {
	go func() {
		pool.registerWorker()
		for {
			// Add ourselves into the worker queue.
			select {
			case work := <-*pool.workQueue():
				work()
			case <-w.quitChan:
				// We have been asked to stop.
				pool.unregisterWorker()
				close(w.quitChan)
				return
			default:
				if w.temp || pool.IsDisposing() {
					w.stop(pool)
				}
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *NaiveWorker) stop(pool WorkerPool) {
	w.shutdownOnce.Do(func() {
		w.quitChan <- true
		fmt.Printf("worker %d stopping\n", w.id)
	})
}
