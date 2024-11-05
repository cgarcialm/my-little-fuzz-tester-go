module github.com/cgarcialm/my-little-fuzz-tester-go

go 1.23
toolchain go1.23.2

// Core TPM library for low-level TPM operations
require github.com/google/go-tpm v0.9.2-0.20240920144513-364d5f2f78b9

require (
    // Utilities for specific TPM-related operations (indirect)
    github.com/google/go-configfs-tsm v0.3.2 // indirect

    // High-level tools for common TPM tasks like key management (indirect)
    github.com/google/go-tpm-tools v0.4.4 // indirect

    // Low-level system calls needed by TPM libraries (indirect)
    golang.org/x/sys v0.15.0 // indirect
)
replace github.com/google/go-tpm => ./internal/github.com/google/go-tpm
