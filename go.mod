module github.com/cgarcialm/my-little-fuzz-tester-go

go 1.23

toolchain go1.23.2

// Core TPM library for low-level TPM operations
require github.com/google/go-tpm v0.9.2-0.20240920144513-364d5f2f78b9

// Low-level system calls needed by TPM libraries (indirect)
require golang.org/x/sys v0.15.0 // indirect
