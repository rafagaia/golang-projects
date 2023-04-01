package main

import (
	"fmt"
	"sync"
)

var sharedValue int

func printEvenNumbers(max int, wg *sync.WaitGroup, priority chan struct{}) {
	defer wg.Done()

	for i := 1; i <= max; i++ {
		select {
		case <-priority:
			if sharedValue%2 == 0 {
				fmt.Println("Even:", sharedValue)
				sharedValue++
			}
			priority <- struct{}{}
		default:
		}
	}
}

func printOddNumbers(max int, wg *sync.WaitGroup, priority chan struct{}) {
	defer wg.Done()

	for i := 1; i <= max; i++ {
		<-priority
		if sharedValue%2 != 0 {
			fmt.Println("Odd:", sharedValue)
			sharedValue++
		}
		priority <- struct{}{}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	priority := make(chan struct{}, 1)
	priority <- struct{}{}

	go printEvenNumbers(10, &wg, priority)
	go printOddNumbers(10, &wg, priority)

	wg.Wait()
}
