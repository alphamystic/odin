package tunnel

import (
	"fmt"
	"net"
)

func RunTCPTunnel() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	targetConn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer targetConn.Close()

	go func() {
		_, err := io.Copy(targetConn, conn)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	_, err = io.Copy(conn, targetConn)
	if err != nil {
		fmt.Println(err)
		return
	}
}
