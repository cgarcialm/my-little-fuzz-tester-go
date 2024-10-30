package main

import (
	"fmt"
	"log"
	"os/exec"
)

func runCommand(name string, args ...string) {
	// Execute the given command with arguments
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command %s failed: %v\nOutput: %s", name, err, string(output))
	}
	fmt.Printf("Command %s output:\n%s\n", name, string(output))
}

func main() {
	// Step 1: Initialize the TPM
	fmt.Println("Initializing the TPM...")
	runCommand("tpm2_startup", "-T", "mssim:host=tpm-simulator,port=2321", "-c")

	// Step 2: Request 16 random bytes from the TPM
	fmt.Println("Requesting 16 random bytes...")
	runCommand("tpm2_getrandom", "16", "-T", "mssim:host=tpm-simulator,port=2321")

	// Step 3: Finalize the TPM session by shutting it down
	fmt.Println("Finalizing the TPM...")
	runCommand("tpm2_shutdown", "-T", "mssim:host=tpm-simulator,port=2321", "-c")
}
