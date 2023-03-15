package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// print a welcome message
	intro()

	/*
	* create a channel to indicate when user wants to quit
	* 	(allows to run things concurrently in the background)
	 */
	doneChan := make(chan bool)

	// start a goroutine to read user input and run program (in background)
	go readUserInput(doneChan)

	/*
	* block until the doneChan gets a value
	*  just listen to doneChan until we get something
	 */
	<-doneChan

	// close the channel
	close(doneChan)

	// exit message
	fmt.Println("...Exiting Program!")
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()

	// check to see if user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	// try to convert what user typed into an int
	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter an Integer!", false
	}

	// ignore first return value, but keep the message
	_, msg := isPrime(numToCheck)

	return msg, false
}

// doesn't return anything because it'll be run as a goroutine
func readUserInput(doneChan chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}
		fmt.Println(res)
		prompt()
	}
}

func intro() {
	fmt.Println("Is it Prime?")
	fmt.Println("-------------")
	fmt.Println("Enter an Integer number, and we'll tell you if it is prime or not.\nEnter 'q' to quit.")
	prompt()
}

func prompt() {
	fmt.Print("-> ") //not Println, because we want cursor to stay on same line
}

func isPrime(num int) (bool, string) {
	// 0 and 1 not prime
	if num == 0 || num == 1 {
		return false, fmt.Sprintf("%d is not prime.", num)
	}
	if num < 0 {
		return false, fmt.Sprintf("Negative numbers (%d) are not prime.", num)
	}

	for i := 2; i <= num/2; i++ {
		if num%i == 0 {
			return false, fmt.Sprintf("%d is not prime. Divisible by %d.", num, i)
		}
	}
	return true, fmt.Sprintf("%d is a prime number.", num)
}
