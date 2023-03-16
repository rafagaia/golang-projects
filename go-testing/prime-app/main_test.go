package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_beta_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number."},
		{"not prime", 8, false, "8 is not prime. Divisible by 2."},
		{"zero", 0, false, "0 is not prime."},
		{"one", 1, false, "1 is not prime."},
		{"negative number", -3, false, "Negative numbers (-3) are not prime."},
	}

	for _, entry := range primeTests {
		result, msg := isPrime(entry.testNum)
		if entry.expected && !result {
			t.Errorf("%s: expected true, got false", entry.name)
		}
		if !entry.expected && result {
			t.Errorf("%s: expected false, got true", entry.name)
		}
		if entry.msg != msg {
			t.Errorf("%s: expected %s, got %s", entry.name, entry.msg, msg)
		}
	}
}

// Test output that's written to the console
func Test_beta_prompt(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe (_ is for err. Ignoring that)
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	prompt()

	//close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt function from our read pipe
	out, _ := io.ReadAll(r)

	// perform our test (cast slice of bytes from ReadAll into string)
	if string(out) != "-> " {
		t.Errorf("Incorrect prompt: expected '-> ' but got '%s'.", string(out))
	}
}

func Test_intro(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe (_ is for err. Ignoring that)
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	intro()

	//close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt function from our read pipe
	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Enter an Integer number") {
		t.Errorf("Intro text not correct; got '%s'", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: "Please enter an Integer!"},
		{name: "zero", input: "0", expected: "0 is not prime."},
		{name: "one", input: "1", expected: "1 is not prime."},
		{name: "two", input: "2", expected: "2 is a prime number."},
		{name: "seven", input: "7", expected: "7 is a prime number."},
		{name: "negative", input: "-5", expected: "Negative numbers (-5) are not prime."},
		{name: "typed_number", input: "five", expected: "Please enter an Integer!"},
		{name: "decimal", input: "2.3", expected: "Please enter an Integer!"},
		{name: "quit", input: "q", expected: ""},
		{name: "QUIT", input: "Q", expected: ""},
		{name: "special_chars", input: "$#@", expected: "Please enter an Integer!"},
	}

	// don't need the index, just the entry
	for _, e := range tests {
		input := strings.NewReader(e.input)
		reader := bufio.NewScanner(input)
		res, _ := checkNumbers(reader)

		if !strings.EqualFold(res, e.expected) {
			t.Errorf("%s: expected %s, but got %s", e.name, e.expected, res)
		}
	}
}

func Test_readUserInput(t *testing.T) {
	// need a channel and instance to an io.Reader
	doneChan := make(chan bool)

	// create a reference to a bytes.Buffer
	//		has required method that satisfies it being a type io.Reader - just need its Read method
	var stdin bytes.Buffer

	stdin.Write([]byte("1\nq\n"))

	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
