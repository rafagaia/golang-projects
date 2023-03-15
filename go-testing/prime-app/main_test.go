package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
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
func Test_prompt(t *testing.T) {
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
