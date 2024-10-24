// File: src/test_orchestrator/test_orchestrator.go

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	// Adjust the import paths to match your module path
	"github.com/cgarcialm/my-little-fuzz-tester-go/src/fuzzer"
	"github.com/cgarcialm/my-little-fuzz-tester-go/src/string_processor"
)

func main() {
	// Read the number of test groups from an environment variable
	numTestGroupsStr := os.Getenv("NUM_TEST_GROUPS")
	if numTestGroupsStr == "" {
		numTestGroupsStr = "1"
	}
	numTestGroups, err := strconv.Atoi(numTestGroupsStr)
	if err != nil {
		fmt.Printf("Invalid NUM_TEST_GROUPS: %v\n", err)
		os.Exit(1)
	}

	// Prepare a directory for logs
	logDir := "/app/logs" // Adjust path if necessary
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Create a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Run tests for each group in parallel
	for i := 0; i < numTestGroups; i++ {
		groupID := i

		wg.Add(1)
		go func(groupID int) {
			defer wg.Done()

			// Create a log file for this group
			logFilePath := filepath.Join(logDir, fmt.Sprintf("group_%d.log", groupID))
			logFile, err := os.Create(logFilePath)
			if err != nil {
				fmt.Printf("Failed to create log file for group %d: %v\n", groupID, err)
				return
			}
			defer logFile.Close()

			// Create a logger
			logger := log.New(logFile, "", log.LstdFlags)

			// Initialize the fuzzer with the string processor function
			f := fuzzer.NewFuzzer(string_processor.ProcessString)
			f.Timeout = 5 * time.Second

			// Run a fixed number of tests
			for j := 0; j < 10; j++ {
				result, err := f.Fuzz()
				if err != nil {
					logger.Printf("Test %d: Error - %v", j+1, err)
				} else {
					logger.Printf("Test %d: Success - %s", j+1, result)
				}
			}

			logger.Printf("Completed tests for group %d", groupID)
		}(groupID)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All test groups have completed.")
}
