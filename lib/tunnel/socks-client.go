package main

import (
	"fmt"
	"net"
	"os"
)

func SClient() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./socks5_client <proxy_address> <target_address>")
		return
	}

	proxyAddress := os.Args[1]
	targetAddress := os.Args[2]

	// Connect to the SOCKS5 proxy server
	proxyConn, err := net.Dial("tcp", proxyAddress)
	if err != nil {
		fmt.Println("Failed to connect to the proxy server:", err)
		return
	}
	defer proxyConn.Close()

	// Send a SOCKS5 version and method selection request
	_, err = proxyConn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		fmt.Println("Failed to send the version/method request:", err)
		return
	}

	// Read the response from the proxy server
	response := make([]byte, 2)
	_, err = proxyConn.Read(response)
	if err != nil {
		fmt.Println("Failed to read the version/method response:", err)
		return
	}

	// Check if the server accepted the method (0x00)
	if response[1] != 0x00 {
		fmt.Println("Proxy server rejected the method.")
		return
	}

	// Create a connection request to the target address
	targetAddr, err := net.ResolveTCPAddr("tcp", targetAddress)
	if err != nil {
		fmt.Println("Failed to resolve target address:", err)
		return
	}

	// Prepare the connection request
	request := make([]byte, 10)
	request[0] = 0x05 // SOCKS5 version
	request[1] = 0x01 // Connect method
	request[3] = 0x01 // IPv4 address type

	// Copy the target address and port to the request
	copy(request[4:8], targetAddr.IP.To4())
	request[8] = byte(targetAddr.Port >> 8)
	request[9] = byte(targetAddr.Port)

	// Send the connection request to the proxy server
	_, err = proxyConn.Write(request)
	if err != nil {
		fmt.Println("Failed to send the connection request:", err)
		return
	}

	// Read the response from the proxy server
	response = make([]byte, 10)
	_, err = proxyConn.Read(response)
	if err != nil {
		fmt.Println("Failed to read the connection response:", err)
		return
	}

	// Check if the connection was established (0x00)
	if response[1] != 0x00 {
		fmt.Println("Proxy server rejected the connection request.")
		return
	}

	// At this point, the connection to the target server is established through the proxy server.
	// You can now send and receive data through the proxy.

	// Example: Send a simple message to the target server
	message := "Hello, SOCKS5 proxy!"
	_, err = proxyConn.Write([]byte(message))
	if err != nil {
		fmt.Println("Failed to send data to the target server:", err)
		return
	}

	// Read the response from the target server
	buffer := make([]byte, 1024)
	n, err := proxyConn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read data from the target server:", err)
		return
	}

	fmt.Printf("Response from the target server: %s\n", string(buffer[:n]))
}
