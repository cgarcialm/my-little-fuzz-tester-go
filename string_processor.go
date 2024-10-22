package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// processString processes the input string and returns an error if the string is too long
func processString(input string) (string, error) {
	// simulate varying times of function
	delay := rand.Intn(6) // Random delay between 0 and 5 seconds
	time.Sleep(time.Duration(delay) * time.Second)

	// Strings longer than 10 characters are invalid
	if len(input) > 10 {
		return "", errors.New("input string is too long")
	}

	// If the string is valid and processing completes in time, return the processed result
	return fmt.Sprintf("Processed: %s", input), nil
}
