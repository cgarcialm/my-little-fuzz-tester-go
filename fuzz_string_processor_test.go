package main

import (
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport/simulator"
)

const (
	ErrInputTooLong = "input string is too long"
)

// func FuzzProcessString(f *testing.F) {
// 	// Seed the fuzzer with initial inputs to guide the fuzzing process
// 	f.Add("Hi")
// 	f.Add("thisisalongstringfortesting")
// 	f.Add("")
// 	f.Add("exactten10")
// 	f.Add("valid@input!")

// 	// Define the fuzz function to test processString
// 	f.Fuzz(func(t *testing.T, input string) {
// 		result, err := processString(input)
// 		if len(input) > 10 {
// 			if err == nil {
// 				t.Errorf("Expected error for input %s, got nil", input)
// 			} else if err.Error() != ErrInputTooLong {
// 				t.Errorf("For input '%s' with length %d, expected error '%s', but got '%s'", input, len(input), ErrInputTooLong, err.Error())
// 			}
// 		} else {
// 			if err != nil {
// 				t.Errorf("For valid input '%s' with length %d, unexpected error: %v", input, len(input), err)
// 			} else {
// 				expected := fmt.Sprintf("Processed: %s", input)
// 				if result != expected {
// 					t.Errorf("Unexpected result for input %s: %s, expected: %s", input, result, expected)
// 				}
// 			}
// 		}
// 	})
// }

// func FuzzGetRandom(f *testing.F) {
// 	// Seed corpus: start with requesting 16 bytes
// 	f.Add(uint16(16))

// 	f.Fuzz(func(t *testing.T, bytesRequested uint16) {
// 		// Open the TPM simulator
// 		tpm, err := simulator.OpenSimulator()
// 		if err != nil {
// 			t.Errorf("Could not open TPM simulator: %v", err)
// 			return
// 		}
// 		defer tpm.Close()

// 		// Validate and constrain the bytesRequested input
// 		if bytesRequested == 0 || bytesRequested > 64 {
// 			return // Skip invalid or excessively large requests
// 		}

// 		// Create the GetRandom command
// 		grc := tpm2.GetRandom{
// 			BytesRequested: bytesRequested,
// 		}

// 		// Execute the command and retrieve the response
// 		response, err := grc.Execute(tpm)
// 		if err != nil {
// 			t.Errorf("GetRandom failed: %v", err)
// 			return
// 		}

// 		// Access the random bytes from the response
// 		randomBytes := response.RandomBytes

// 		// Verify the length of the returned bytes
// 		if len(randomBytes.Buffer) != int(bytesRequested) {
// 			t.Errorf("Expected %d bytes, got %d", bytesRequested, len(randomBytes.Buffer))
// 			return
// 		}

//			// Print the random bytes as a hex string
//			fmt.Printf("Random Bytes: %x\n", randomBytes.Buffer)
//		})
//	}
func FuzzTPMGetRandom(f *testing.F) {
	// Seed with initial input for the number of bytes requested
	f.Add(16)

	f.Fuzz(func(t *testing.T, numBytes int) {
		// if numBytes == 20 { // Force an interesting input by defining a condition
		// 	t.Fatalf("Forced failure for numBytes == 20")
		// }
		f.Add(16)

		if numBytes < 1 || numBytes > 64 {
			t.Skip("Number of bytes out of range")
		}

		tpm, err := simulator.OpenSimulator()
		if err != nil {
			t.Fatalf("Failed to connect to TPM simulator: %v", err)
		}
		defer tpm.Close()

		// Run the GetRandom command
		grc := tpm2.GetRandom{
			BytesRequested: uint16(numBytes),
		}
		_, err = grc.Execute(tpm)
		if err != nil {
			t.Logf("GetRandom failed: %v", err)
		}
	})
}
