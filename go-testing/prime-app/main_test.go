package main

import "testing"

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number."},
		{"not prime", 8, false, "8 is not prime. Divisible by 2"},
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
