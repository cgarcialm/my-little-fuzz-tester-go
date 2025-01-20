package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// TPM represents a connection to a TPM simulator.
type TPM struct {
	transport *net.TCPConn
}

// Send implements the TPM interface.
func (t *TPM) Send(input []byte) ([]byte, error) {
	fmt.Println("Sending bytes to TPM.")
	_, err := t.transport.Write(input)
	if err != nil {
		fmt.Println("Write to server failed:", err.Error())
		return nil, err
	}

	reply := make([]byte, 1024)

	// Set a timeout for the write and read operations
	timeout := 5 * time.Second // Adjust the timeout duration as needed

	// Set the deadline for reading
	err = t.transport.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		fmt.Println("SetReadDeadline failed:", err.Error())
		return nil, err
	}

	fmt.Println("Reading response from TPM.")
	n, err := t.transport.Read(reply)
	if err != nil {
		fmt.Println("Read from server failed:", err.Error())
		return nil, err
	}

	// // Check how much data was actually read
	fmt.Printf("Read %d bytes\n", n)

	return reply, nil
}

// Close implements the TPM interface.
func (t *TPM) Close() error {
	return t.transport.Close()
}

func main() {
	var err error
	servAddr := "tpm-simulator:2321"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}

	// var thetpm transport.TPMCloser
	thetpm := &TPM{
		transport: conn,
	}
	defer thetpm.Close()

	var output []byte
	var message []byte

	message = []byte{
		0, 0, 0, 0x0F, // TPM_REMOTE_HANDSHAKE
		0, 0, 0, 0x04, // Client Version
	}

	fmt.Printf("Sending handshake. Message: %x\n", message)
	output, err = thetpm.Send(message)
	if err != nil {
		fmt.Println("Remote handshake failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Remote Handshake successful. Output: %x ... %x\n", output[:32], output[len(output)-32:])

	message = []byte{
		0, 0, 0, 0x08, // TPM_SEND_COMMAND
		0, 0, 0, 0x0C, // length
		0, 0, 0x01, 0x7B, // TPM_CC_GetRandom
		0, 0x08, // Command parameter - num random bytes to generate
	}

	output, err = thetpm.Send(message)
	fmt.Printf("Sending get random. Message: %x\n", message)
	if err != nil {
		fmt.Println("Get random command failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Get random successful. Output: %x ... %x\n", output[:32], output[len(output)-32:])

	message = []byte{
		0, 0, 0, 0x15, // TPM_STOP
	}
	fmt.Printf("Sending stop. Message: %x\n", message)
	output, err = thetpm.Send(message)
	if err != nil {
		fmt.Println("Stop failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Stop successful. Output: %x ... %x\n", output[:32], output[len(output)-32:])

	// Goal is to have this implementation for each command
	// so that we can call "Execute" and have the TCP
	// communication happen.

	// grc := GetRandom{
	// 	BytesRequested: 16,
	// }

	// var out *GetRandomResponse
	// if out, err = grc.Execute(thetpm); err != nil {
	// 	fmt.Printf("GetRandom failed: %v", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Output: %x", out.RandomBytes)
}
