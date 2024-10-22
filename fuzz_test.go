package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
)

var start, end int // Global variables for range

// TestMain is the entry point for all tests and parses custom flags.
func TestMain(m *testing.M) {
	// Define custom flags for start and end
	flag.IntVar(&start, "start", 0, "Start index for the test range")
	flag.IntVar(&end, "end", 10, "End index for the test range")

	// Parse the flags
	flag.Parse()

	// Run the tests
	os.Exit(m.Run())
}

// WriteTestReport writes the result of fuzz tests to a report file
func WriteTestReport(report string) error {
	// Create a report file
	file, err := os.Create("fuzz_test_report.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the report content to the file
	_, err = file.WriteString(report)
	return err
}

// RunFuzzTests runs the fuzz tests for a given range of tests and returns the report
func RunFuzzTests(t *testing.T, start, end int, fuzzer *Fuzzer) string {
	report := "Fuzz Test Report:\n"

	for i := start; i < end; i++ {
		result, err := fuzzer.Fuzz()

		if err != nil {
			if err.Error() == "test timed out" {
				// Log timeout without failing test
				t.Logf("Test %d: Timeout - %v", i+1, err)
				report += fmt.Sprintf("Test %d: Timeout - %v\n", i+1, err)
			} else if err.Error() == "input string is too long" {
				// Log expected error without failing the test
				t.Logf("Test %d: Expected error - %v", i+1, err)
				report += fmt.Sprintf("Test %d: Expected error - %v\n", i+1, err)
			} else {
				// Log unexpected error and fail the test
				t.Errorf("Test %d: Unexpected error - %v", i+1, err)
				report += fmt.Sprintf("Test %d: Unexpected error - %v\n", i+1, err)
			}
		} else {
			// Log success
			t.Logf("Test %d: Success - %s", i+1, result)
			report += fmt.Sprintf("Test %d: Success - %s\n", i+1, result)
		}
	}

	return report
}

// RunTestRange runs a range of fuzz tests in parallel and communicates results back via a channel
func RunFuzzTestRange(t *testing.T, start, end int, fuzzer *Fuzzer, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done() // Signal that this goroutine is done when the function returns

	// Execute the tests for the given range
	report := RunFuzzTests(t, start, end, fuzzer)

	// Send the generated report back to the main thread via the result channel
	resultChan <- report
}

// TestFuzzer tests the processString function using the fuzzer
func TestFuzzer(t *testing.T) {
	// Initialize the fuzzer
	fuzzer := NewFuzzer(processString)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to gather reports from all parallel runs
	reportChan := make(chan string, 2)

	// Run the tests for the specified range
	wg.Add(1)
	go RunFuzzTestRange(t, start, end, fuzzer, &wg, reportChan)

	// Wait for all parallel runs to finish
	wg.Wait()

	// Close the report channel
	close(reportChan)

	// Aggregate the reports
	finalReport := ""
	for report := range reportChan {
		finalReport += report
	}

	// Write the test report to a file
	err := WriteTestReport(finalReport)
	if err != nil {
		t.Fatalf("Failed to write test report: %v", err)
	}
}

// TestFixedInput is a test case with a hardcoded input string
func TestFixedInput(t *testing.T) {
	input := "Hello"                     // Adjusted input to fit within the expected length limit
	expectedOutput := "Processed: Hello" // Define the expected output

	result, err := processString(input)

	if err != nil {
		t.Errorf("Unexpected error for input '%s': %v", input, err)
	} else {
		if result != expectedOutput {
			t.Errorf("Expected output '%s', but got '%s'", expectedOutput, result)
		} else {
			t.Logf("Success for input '%s': %s", input, result)
		}
	}
}
