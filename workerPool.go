package workerpool

import (
	"log"
	"sync"
	"sync/atomic"
)

// WorkerPool is a simple thread safe worker pool
type WorkerPool interface {
	Submit(work func())
	Dispose()
	IsDisposing() bool
	IsDisposed() bool
	workQueue() *chan func()
	registerWorker()
	unregisterWorker()
}

// NaiveWorkerPool provides a naive implementation of the worker pool interface
type NaiveWorkerPool struct {
	initialSize  uint64
	maxWorker    uint64
	workerCount  uint64
	work         chan func()
	wg           sync.WaitGroup
	disposing    int32
	disposed     int32
	shutdownOnce sync.Once
}

// Create a naive worker pool with the specified initial and max size
func Create(initialSize uint64, maxSize uint64) *NaiveWorkerPool {
	log.Printf("Creating worker pool with size %d and max size %d\n", initialSize, maxSize)
	nwp := &NaiveWorkerPool{
		initialSize:  initialSize,
		maxWorker:    maxSize,
		work:         make(chan func(), 500),
		shutdownOnce: sync.Once{},
	}
	for i := int64(0); uint64(i) < initialSize; i++ {
		nwp.spawnWorker(false)
	}

	return nwp
}

// Submit work to the worker pool. Work will not be submitted if the pool is shutting down or shut down.
func (nwp *NaiveWorkerPool) Submit(workerFunc func()) {
	if nwp.IsDisposed() {
		log.Panicln("Cannot submit work to a pool which has been shut down!")
	} else if nwp.IsDisposing() {
		log.Print("Cannot submit work to a pool which is shutting down!")
	} else {
		nwp.work <- workerFunc
		nwp.spawnNewWorkerIfNeeded()
	}
}

// Dispose shuts down the worker pool, but waits for all work to finish
func (nwp *NaiveWorkerPool) Dispose() {
	if !nwp.IsDisposed() && !nwp.IsDisposing() {
		nwp.shutdownOnce.Do(nwp.internalShutdown)
	} else if nwp.IsDisposed() {
		log.Println("Worker Pool already disposed.")
	} else if nwp.IsDisposing() {
		log.Println("Worker Pool already shutting down.")
	}
}

// IsDisposed checks if the pool has been disposed
func (nwp *NaiveWorkerPool) IsDisposed() bool {
	return atomic.LoadInt32(&nwp.disposed) == 1
}

//IsDisposing checks if the pool is currently disposing/shutting down
func (nwp *NaiveWorkerPool) IsDisposing() bool {
	return atomic.LoadInt32(&nwp.disposing) == 1
}

func (nwp *NaiveWorkerPool) registerWorker() {
	nwp.wg.Add(1)
}

func (nwp *NaiveWorkerPool) unregisterWorker() {
	atomic.AddUint64(&nwp.workerCount, ^uint64(0))
	nwp.wg.Done()
}

func (nwp *NaiveWorkerPool) internalShutdown() {
	atomic.StoreInt32(&nwp.disposing, 1)
	log.Println("Shutting down worker pool...")
	close(nwp.work)
	nwp.wg.Wait()
	log.Println("Worker pool shut down.")
	atomic.StoreInt32(&nwp.disposed, 1)
}

func (nwp *NaiveWorkerPool) workQueue() (result *chan func()) {
	result = &nwp.work
	return
}

func (nwp *NaiveWorkerPool) spawnNewWorkerIfNeeded() {
	if len(nwp.work) > 10 && atomic.LoadUint64(&nwp.workerCount) < nwp.maxWorker {
		nwp.spawnWorker(true)
	}
}

func (nwp *NaiveWorkerPool) setDisposing() bool {
	return atomic.SwapInt32(&nwp.disposing, 1) == 0
}

func (nwp *NaiveWorkerPool) spawnWorker(temp bool) {
	id := atomic.AddUint64(&nwp.workerCount, 1)
	log.Printf("Spawning worker %d is temp: %t\n", id, temp)

	worker := &NaiveWorker{
		temp: temp,
		id:   id,
	}
	worker.start(nwp)
}
