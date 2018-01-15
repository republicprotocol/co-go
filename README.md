# Concurrency

[![Build Status](https://travis-ci.org/republicprotocol/go-async.svg?branch=master)](https://travis-ci.org/republicprotocol/go-async)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-async/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-async?branch=master)

The Concurrency library is a Go implementation of high level concurrent features. It provides a simple API for common task parallel and data parallel constructs. Using goroutines, the parallelism provided is actually a form of concurrency since goroutines are not guaranteed to run strict simultaneity.

## For all

The `ForAll` loop is a data parallel loop that distributes iterations evenly across several goroutines. It will launch one goroutine per CPU, and can be used on arrays, maps, and slices.

```
xs := []int{1,2,3,4,5,6,7,8,10}
do.ForAll(xs, func(i int) {
    xs[i] *= 2
})
```

It is the responsibility of the programmer to ensure that the function being used is safe for concurrent environments. A simple way of ensuring this is checking that the function will never mutate any object other than the object accessible using the `i` index. You can also use the go tools to check for race conditions during testing.

## Tests

To run the test suite, install Ginkgo.

```
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```
ginkgo -v --trace --cover --coverprofile coverprofile.out
```

## License

The Concurrency library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.