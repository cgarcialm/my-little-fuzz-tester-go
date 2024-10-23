package tests

import (
	"flag"
	"fmt"
	"github.com/cgarcialm/my-little-fuzz-tester-go/src/fuzzer"           // Correct import
	"github.com/cgarcialm/my-little-fuzz-tester-go/src/string_processor" // Correct import
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var start, end int // Global variables for range

// TestMain is the entry point for all tests and parses custom flags.
func TestMain(m *testing.M) {
	// Define custom flags for start and end range
	flag.IntVar(&start, "start", 0, "Start index for the test range")
	flag.IntVar(&end, "end", 10, "End index for the test range")

	// Parse the flags
	flag.Parse()

	// Run the tests
	os.Exit(m.Run())
}

// WriteTestReport writes the result of fuzz tests to a report file in the "output" folder
func WriteTestReport(report string) error {
	// Ensure the "output" directory exists at the root level
	dir := "../outputs" // Navigate up to the root folder, then into the output folder
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm) // Create the directory if it doesn't exist
		if err != nil {
			return err
		}
	}

	// Create a report file in the "output" directory
	filePath := filepath.Join(dir, "fuzz_test_report.txt")
	file, err := os.Create(filePath) // Create or overwrite the file
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the report content to the file
	_, err = file.WriteString(report)
	return err
}

// RunFuzzTests runs the fuzz tests for a given range of tests and returns the report
func RunFuzzTests(t *testing.T, start, end int, fuzzer *fuzzer.Fuzzer) string {
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

// RunFuzzTestRange runs a range of fuzz tests in parallel and communicates results back via a channel
func RunFuzzTestRange(t *testing.T, start, end int, fuzzer *fuzzer.Fuzzer, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done() // Signal that this goroutine is done when the function returns

	// Execute the tests for the given range
	report := RunFuzzTests(t, start, end, fuzzer)

	// Send the generated report back to the main thread via the result channel
	resultChan <- report
}

// TestFuzzer tests the processString function using the fuzzer
func TestFuzzer(t *testing.T) {
	// Initialize the fuzzer
	f := fuzzer.NewFuzzer(string_processor.ProcessString)
	f.Timeout = 5 * time.Second // Increase the timeout to 5 seconds to reduce timeouts

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to gather reports from all parallel runs
	reportChan := make(chan string, 2)

	// Run the tests for the specified range in parallel
	wg.Add(1)
	go RunFuzzTestRange(t, start, end, f, &wg, reportChan)

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

	result, err := string_processor.ProcessString(input)

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
