package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	tpmAddress = flag.String("tpm-address", "tpm-simulator:2321", "Address of the TPM simulator (host:port).")
	pcr        = flag.Int("pcr", 0, "PCR index to seal data to (must be within [0, 23]).")
)

func main() {
	flag.Parse()

	if *pcr < 0 || *pcr > 23 {
		log.Fatalf("Invalid PCR index: %d. Must be within [0, 23].", *pcr)
	}

	// Connect to the TPM simulator over TCP
	conn, err := net.Dial("tcp", *tpmAddress)
	if err != nil {
		log.Fatalf("Failed to connect to TPM simulator: %v", err)
	}
	fmt.Println("Connection to TPM simulator established!")
	defer conn.Close()
}
