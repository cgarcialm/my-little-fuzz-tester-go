package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Fuzzer struct that holds a function to be fuzzed and a timeout
type Fuzzer struct {
	testFunction func(string) (string, error) // The function to be fuzzed
	timeout      time.Duration                // Timeout duration for each test
}

// NewFuzzer creates a new fuzzer instance with a timeout
func NewFuzzer(testFunction func(string) (string, error)) *Fuzzer {
	return &Fuzzer{
		testFunction: testFunction,
		timeout:      3 * time.Second, // 3-second timeout
	}
}

// randomString generates a random string of a given length
func (f *Fuzzer) randomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s := make([]byte, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Fuzz generates random inputs, tests the function with a timeout, and returns the result and any error
func (f *Fuzzer) Fuzz() (string, error) {
	input := f.randomString(rand.Intn(20) + 1) // Generate a random string of length 1 to 20
	resultChan := make(chan string)
	errorChan := make(chan error)

	// Run the test function in a separate goroutine
	go func() {
		result, err := f.testFunction(input) // Call the function being fuzzed
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()

	// Use select to handle either a result, an error, or a timeout
	select {
	case result := <-resultChan:
		fmt.Printf("Success with input '%s': %s\n", input, result)
		return result, nil
	case err := <-errorChan:
		fmt.Printf("Error with input '%s': %v\n", input, err)
		return "", err
	case <-time.After(f.timeout):
		fmt.Printf("Test timed out after %v seconds for input: %s\n", f.timeout.Seconds(), input)
		return "", errors.New("test timed out")
	}
}
