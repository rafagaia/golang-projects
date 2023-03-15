package main

import (
	"fmt"
)

func main() {
	// print welcome message
	/*intro()

	// create a channel to indicate when the user wants to quit
	doneChan := make(chan bool)

	// start a goroutine to read user input and run program
	go readUserInput(os.Stdin, doneChan)

	// block until the doneChan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	// exit message
	fmt.Println("Exiting.")
	*/

	n := 0

	_, msg := isPrime(n)
	fmt.Println(msg)
}

func isPrime(num int) (bool, string) {
	// 0 and 1 not prime
	if num == 0 || num == 1 {
		return false, fmt.Sprintf("%d is not prime.", num)
	}
	if num < 0 {
		return false, fmt.Sprintf("%d negative numbers not prime", num)
	}

	for i := 2; i <= num/2; i++ {
		if num%i == 0 {
			return false, fmt.Sprintf("%d is not prime. Divisible by %d", num, i)
		}
	}
	return true, fmt.Sprintf("%d is a prime number.", num)
}

/*func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)

	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}

		fmt.Println(res)
		// ....
	}

	// ....
}*/

// ....
