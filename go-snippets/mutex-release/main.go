package main

/*
	Protecting access to shared state: When multiple goroutines access
	shared data, a mutex can be used to ensure that only one goroutine
	can modify the data at a time, preventing race conditions.

	a mutex (short for "mutual exclusion") is a synchronization
	primitive provided by the Go sync package. It is used to ensure
	that only one goroutine can access a shared resource or a critical
	section of code at a time. Mutexes help prevent race conditions and
	ensure that concurrent access to shared data is safe.
**/

import "sync"

var mu sync.Mutex

/*
data is a shared resource.
maps are not safe for concurrent use in Go.
thus, it is necessary to use a mutex to protect access to
the shared resource
*/
var data map[string]string

/*
set value to shared map in a way that ensures that concurrent
calls to setValue will not create race conditions when
accessing the shared map
*/
func setValue(key, value string) {
	// acquire mutex lock
	mu.Lock()
	// ensure lock is released when function exits
	defer mu.Unlock()

	data[key] = value
}

func main() {
	go setValue("1", "new value")
	go setValue("1", "updated value")
	go setValue("1", "updated value once more")
}
