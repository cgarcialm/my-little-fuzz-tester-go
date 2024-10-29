package main

import (
	"fmt"
	"testing"
)

const (
	ErrInputTooLong = "input string is too long"
)

func FuzzProcessString(f *testing.F) {
	// Seed the fuzzer with initial inputs to guide the fuzzing process
	f.Add("Hi")
	f.Add("thisisalongstringfortesting")
	f.Add("")
	f.Add("exactten10")
	f.Add("valid@input!")

	// Define the fuzz function to test processString
	f.Fuzz(func(t *testing.T, input string) {
		result, err := processString(input)
		if len(input) > 10 {
			if err == nil {
				t.Errorf("Expected error for input %s, got nil", input)
			} else if err.Error() != ErrInputTooLong {
				t.Errorf("For input '%s' with length %d, expected error '%s', but got '%s'", input, len(input), ErrInputTooLong, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("For valid input '%s' with length %d, unexpected error: %v", input, len(input), err)
			} else {
				expected := fmt.Sprintf("Processed: %s", input)
				if result != expected {
					t.Errorf("Unexpected result for input %s: %s, expected: %s", input, result, expected)
				}
			}
		}
	})
}
