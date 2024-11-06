package main

import (
	"fmt"
	"testing"
	"time"
)

// Test that valid strings are processed correctly
func TestProcessString_ValidInput(t *testing.T) {
	input := "Testing"
	expectedOutput := "Processed: Testing"

	// Call processString with a valid input
	result, err := processString(input)

	if err != nil {
		t.Errorf("Unexpected error for valid input '%s': %v", input, err)
	}

	// Ensure the result matches the expected output
	if result != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, result)
	} else {
		fmt.Printf("TestProcessString_ValidInput passed\n")
	}
}

// Test that strings longer than 10 characters will result in an error
func TestProcessString_TooLongInput(t *testing.T) {
	input := "ThisStringIsTooLong"

	// Call processString with an invalid (too long) input
	_, err := processString(input)

	// Ensure the error matches the expected "too long" error
	if err == nil || err.Error() != "input string is too long" {
		t.Errorf("Expected 'input string is too long' error, but got '%v'", err)
	} else {
		fmt.Printf("TestProcessString_TooLongInput passed\n")
	}
}

// Test that an empty string still returns
func TestProcessString_EmptyString(t *testing.T) {
	input := ""
	expectedOutput := "Processed: "

	// Call processString with an empty string
	result, err := processString(input)

	if err != nil {
		t.Errorf("Unexpected error for empty input: %v", err)
	}

	// Ensure the result matches the expected output
	if result != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, result)
	} else {
		fmt.Printf("TestProcessString_EmptyString passed\n")
	}
}

// Test the simulated delay in processString function, throws an error if is more than 5 seconds.
func TestProcessString_SimulatedDelay(t *testing.T) {
	input := "DelayTest"

	// Measure the time taken to process the string
	startTime := time.Now()
	_, err := processString(input)
	elapsed := time.Since(startTime)

	// Ensure the function finishes within a reasonable range (0-5 seconds)
	if elapsed > 6*time.Second {
		t.Errorf("processString took too long: %v", elapsed)
	}

	if err != nil {
		t.Errorf("Unexpected error for input '%s': %v", input, err)
	} else {
		fmt.Printf("TestProcessString_SimulatedDelay passed\n")
	}
}
