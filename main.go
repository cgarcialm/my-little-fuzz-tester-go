package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/go-tpm/tpm2"
)

var tpmAddress = flag.String("tpm-address", "tpm-simulator:2321", "Address of the TPM simulator (host:port)")

// ConnTransport wraps net.Conn to implement the transport.TPM interface
type ConnTransport struct {
	conn net.Conn
}

// Send sends data to the TPM and receives a response.
func (c *ConnTransport) Send(data []byte) ([]byte, error) {
	if _, err := c.conn.Write(data); err != nil {
		return nil, fmt.Errorf("failed to send data to TPM: %v", err)
	}

	resp := make([]byte, 4096)
	n, err := c.conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from TPM: %v", err)
	}

	return resp[:n], nil
}

func (c *ConnTransport) Close() error {
	return c.conn.Close()
}

func main() {
	flag.Parse()

	// Connect to the TPM simulator over TCP
	conn, err := net.Dial("tcp", *tpmAddress)
	if err != nil {
		log.Fatalf("Failed to connect to TPM simulator: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to TPM simulator!")

	// Wrap the connection in a ConnTransport to meet the transport.TPM interface
	tpmTransport := &ConnTransport{conn: conn}

	// Create a GetRandom command for 16 bytes
	cmd := tpm2.GetRandom{BytesRequested: 16}

	// Execute the command
	resp, err := cmd.Execute(tpmTransport)
	if err != nil {
		log.Fatalf("Failed to execute GetRandom command: %v", err)
	}

	// Print the random bytes from TPM
	fmt.Printf("Random bytes from TPM: %x\n", resp.RandomBytes)
}
