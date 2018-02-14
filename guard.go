package do

import (
	"sync"
)

// A Guard is a conditional variable used for entry into a function. It allows
// a function to wait for a condition to be true before executing. It is not
// safe to use concurrently.
type Guard struct {
	condition     func() bool
	conditionOpen bool
	conditionWait chan struct{}
}

func newGuard(condition func() bool) Guard {
	guard := Guard{
		condition:     condition,
		conditionOpen: condition(),
		conditionWait: make(chan struct{}, 1),
	}
	if guard.conditionOpen {
		guard.conditionWait <- struct{}{}
	}
	return guard
}

func (guard *Guard) open() {
	if guard.conditionOpen {
		return
	}
	guard.conditionOpen = true
	guard.conditionWait <- struct{}{}
}

func (guard *Guard) close() {
	if !guard.conditionOpen {
		return
	}
	guard.conditionOpen = false
	<-guard.conditionWait
}

func (guard *Guard) wait() {
	<-guard.conditionWait
	guard.conditionOpen = false
}

// A GuardedObject uses Guards to provide safe concurrent access to an object.
// The object should use a GuardedObject to create Guards, and each function
// that accesses the object must call the Enter and Exit functions at the
// beginning and end of the function. A Guard can optionally be passed to the
// Enter function, as long as the Guard was created using the same
// GuardedObject. By doing this, the function will no execute until the Guard
// condition is true.
type GuardedObject struct {
	mu     *sync.RWMutex
	guards []Guard
}

func NewGuardedObject() GuardedObject {
	return GuardedObject{
		mu:     new(sync.RWMutex),
		guards: make([]Guard, 0),
	}
}

// Guard returns a Guard for use with this GuardedObject, that waits for the
// condition to be true. The Guard condition must be read-only and make no
// changes to non-local variables.
func (object *GuardedObject) Guard(condition func() bool) *Guard {
	guard := newGuard(condition)
	object.mu = new(sync.RWMutex)
	object.guards = append(object.guards, guard)
	return &object.guards[len(object.guards)-1]
}

// Enter the GuardedObject. Exit must be called after a call to Enter. If a
// Guard is passed to this function it will block until the Guard condition
// is true.
func (object *GuardedObject) Enter(guard *Guard) {
	if guard != nil {
		guard.wait()
	}
	object.mu.Lock()
}

// EnterReadOnly will enter the GuardedObject but only acquire a read lock. Any
// function that uses EnterReadOnly to protect an object must make sure that it
// does not modify the object. ExitReadOnly must be called after a call to
// EnterReadOnly.
func (object *GuardedObject) EnterReadOnly(guard *Guard) {
	if guard != nil {
		guard.wait()
	}
	object.mu.RLock()
}

// Exit the GuardedObject. All Guards attached to the GuardedObject will be
// re-evaluated. This must not be called unless a call to Enter has already
// been made.
func (object *GuardedObject) Exit() {
	object.resolveGuards()
	object.mu.Unlock()
}

// ExitReadOnly is the same as Exit, but it must not be called unless a call to
// EnterReadOnly has already been made.
func (object *GuardedObject) ExitReadOnly() {
	object.resolveGuards()
	object.mu.RUnlock()
}

func (object *GuardedObject) resolveGuards() {
	for i := range object.guards {
		if object.guards[i].condition() {
			object.guards[i].open()
		} else {
			object.guards[i].close()
		}
	}
}
