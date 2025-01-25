package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// TPM represents a connection to a TPM simulator.
type TPM struct {
	transport *net.TCPConn
}

func wrapCommandForTcp(input []byte) []byte {
	// The TPM will always be used to send simple commands
	// so hardcode the command type.
	command := make([]byte, 4)
	binary.BigEndian.PutUint32(command, uint32(SendCommand))
	// Always send locality zero.
	command = append(command, uint8(0))
	// Measure the input and add its length.
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(input)))
	command = append(command, length...)
	// Finally, add the payload.
	return append(command, input...)
}

// Send implements the TPM interface.
func (t *TPM) Send(input []byte) ([]byte, error) {
	fmt.Printf("Marshalled command: %x\n", input)
	command := wrapCommandForTcp(input)
	fmt.Printf("Wrapped command: %x\n", command)

	_, err := t.transport.Write(command)
	if err != nil {
		println("Write to server failed:", err.Error())
		return nil, err
	}

	length_buffer := make([]byte, 4)
	_, err = t.transport.Read(length_buffer)
	if err != nil {
		println("Read from server failed:", err.Error())
		return nil, err
	}

	length := binary.BigEndian.Uint32(length_buffer)
	fmt.Printf("Length is: %d\n", length)
	if length > MaxBufferSize {
		fmt.Printf("Payload is bigger than max buffer size: %d", length)
		return nil, err
	}

	payload_buffer := make([]byte, length)
	_, err = t.transport.Read(payload_buffer)
	if err != nil {
		fmt.Printf("Read from server failed: %v", err.Error())
		return nil, err
	}

	end_command_buffer := make([]byte, 4)
	_, err = t.transport.Read(end_command_buffer)
	if err != nil {
		fmt.Printf("Read from server failed: %v", err.Error())
		return nil, err
	}
	end_command := binary.BigEndian.Uint32(end_command_buffer)
	if end_command != 0 {
		fmt.Printf("Unexpected end command: %d\n", end_command)
	}

	return payload_buffer, nil
}

// Close implements the TPM interface.
func (t *TPM) Close() error {
	return t.transport.Close()
}

type TPMTcpClient struct {
	commandClient  *net.TCPConn
	platformClient *net.TCPConn
}

func (t *TPMTcpClient) Close() error {
	if err := t.commandClient.Close(); err != nil {
		return err
	}
	return t.platformClient.Close()
}

func (t *TPMTcpClient) SendPlatformCommand(command TcpTpmCommand) error {
	command_buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(command_buffer, uint32(command))

	_, err := t.platformClient.Write(command_buffer)
	if err != nil {
		println("Write to server failed:", err.Error())
		return err
	}

	end_command_buffer := make([]byte, 4)
	_, err = t.platformClient.Read(end_command_buffer)
	if err != nil {
		fmt.Printf("Read from server failed: %v", err.Error())
		return err
	}
	end_command := binary.BigEndian.Uint32(end_command_buffer)
	if end_command != 0 {
		fmt.Printf("Unexpected end command: %d\n", end_command)
	}
	return nil
}

func (t *TPMTcpClient) SetUpTpm() error {
	// Need to power on the TPM to send commands
	if err := t.SendPlatformCommand(SignalPowerOn); err != nil {
		return err
	}
	return nil
}

func (t *TPMTcpClient) GetTPM() *TPM {
	return &TPM{
		transport: t.commandClient,
	}
}

func startTcpConnection(servAddr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	return conn, nil
}

func CreateTPMTcpClient(commandAddr, platformAddr string) (*TPMTcpClient, error) {
	commandClient, err := startTcpConnection(commandAddr)
	if err != nil {
		println("command client failed:", err.Error())
		return nil, err
	}
	platformClient, err := startTcpConnection(platformAddr)
	if err != nil {
		println("platform client failed:", err.Error())
		return nil, err
	}

	return &TPMTcpClient{
		commandClient:  commandClient,
		platformClient: platformClient,
	}, nil
}
