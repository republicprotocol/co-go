package do

import (
	"reflect"
	"runtime"
	"sync"
)

// ForAll items in the data set, apply the function. The function accepts the
// index of the item to which is it being applied. One goroutine is launched
// for each CPU, so the given function must be safe to use concurrently.
func ForAll(data interface{}, f func(i int)) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		// Calculate workload size per CPU.
		length := reflect.ValueOf(data).Len()
		numCPUs := runtime.NumCPU()
		numIterationsPerCPU := (length / numCPUs) + 1
		// Apply the function in parallel over the data.
		var wg sync.WaitGroup
		wg.Add(numCPUs)
		for offset := 0; offset < length; offset += numIterationsPerCPU {
			go func(offset int) {
				defer wg.Done()
				for i := offset; i < offset+numIterationsPerCPU && i < length; i++ {
					f(i)
				}
			}(offset)
		}
		wg.Wait()
	}
}

// Option types are returned from Process functions. They are either a value or
// an error. The error should be checked before using the value.
type Option struct {
	Value interface{}
	Err   error
}

// Some returns an Option struct with a value and no error.
func Some(v interface{}) Option {
	return Option{
		Value: v,
	}
}

// Err returns a Return struct with an error and no value.
func Err(err error) Option {
	return Option{
		Err: err,
	}
}

// Process runs the function in a goroutine and writes the return value to a
// channel.
func Process(f func() Option) chan Option {
	ch := make(chan Option)
	go func() {
		ch <- f()
	}()
	return ch
}
