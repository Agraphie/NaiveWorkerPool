NaiveWorkerPool
==========
[![Go Report Card](https://goreportcard.com/badge/github.com/Agraphie/NaiveWorkerPool)](https://goreportcard.com/badge/github.com/Agraphie/NaiveWorkerPool)

Go package for using a simple and dynamically growing worker pool to submit functions to. 

```go
initialWorkerCount := uint64(4)
maxWorker := uint(64)
pool := Create(initialWorkerCount, maxWorker)

pool.Submit(func() {
  // Any work you want to be done. Pool will grow dynamically when too much work is queued up.
})

// Shut down the pool, but wait for all work to be done
pool.Dispose()
```

A pool can grow until it reaches the specified max worker limit. The pool will stay on that level until all work is done. It will
then decrease again until it reaches the initial size. The growing will (naively) starting now when more than 10 functions are queued up.

See [GoDoc](https://godoc.org/github.com/Agraphie/NaiveWorkerPool) for more details.
