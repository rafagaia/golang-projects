package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var sharedValue int

func printEvenNumbers(max int, wg *sync.WaitGroup) {
	// decrements WaitGroup counter by 1
	defer wg.Done()

	for i := 1; i <= max; i++ {
		mu.Lock()
		if sharedValue%2 == 0 {
			fmt.Println("Even:", sharedValue)
			sharedValue++
		}
		mu.Unlock()
	}
}

func printOddNumbers(max int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= max; i++ {
		mu.Lock()
		if sharedValue%2 != 0 {
			fmt.Println("Odd:", sharedValue)
			sharedValue++
		}
		mu.Unlock()
	}
}

func main() {
	/*
	* create synchronization primitive:
	*  group to coordinate completion of multiple goroutines.
	*  helps in waiting for a collection of goroutines to finish executing.
	 */
	var wg sync.WaitGroup
	// increment WaitGroup counter to indicate that 2 goroutines will be launched
	wg.Add(2)

	go printEvenNumbers(10, &wg)
	go printOddNumbers(10, &wg)

	// blocks caller until WaitGroup counter reaches 0
	wg.Wait()
}
