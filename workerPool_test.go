package workerpool

import (
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	initialWorkerCount := uint64(4)
	pool := Create(initialWorkerCount, 10)

	if pool.workerCount != initialWorkerCount {
		t.Errorf("Worker sum was incorrect, got: %d, want: %d.", pool.workerCount, initialWorkerCount)
	}
}

func TestDisposePool(t *testing.T) {
	initialWorkerCount := uint64(4)
	pool := Create(initialWorkerCount, 10)
	pool.Dispose()

	if !pool.IsDisposed() {
		t.Errorf("Expected pool to be disposed, but was not.")
	}
}

func TestDynamicallyGrowing(t *testing.T) {
	initialWorkerCount := uint64(4)
	maxWorker := uint64(10)
	pool := Create(initialWorkerCount, maxWorker)
	for i := 0; i < 100; i++ {
		pool.Submit(func() {
			time.Sleep(time.Second)
		})
	}

	if pool.workerCount < maxWorker {
		t.Errorf("Worker sum was incorrect, got: %d, want: %d.", pool.workerCount, maxWorker)
	}
}

func TestDynamicallyDecreasing(t *testing.T) {
	initialWorkerCount := uint64(4)
	maxWorker := uint64(10)
	pool := Create(initialWorkerCount, maxWorker)

	for i := 0; i < 20; i++ {
		pool.Submit(func() {
			time.Sleep(2 * time.Second)
		})
	}

	if pool.workerCount != maxWorker {
		t.Errorf("Expected pool to be its maximum size, got: %d, want: %d.", pool.workerCount, maxWorker)
	}

	time.Sleep(20 * time.Second)
	if pool.workerCount != initialWorkerCount {
		t.Errorf("Worker sum was incorrect, got: %d, want: %d.", pool.workerCount, initialWorkerCount)
	}
}
