package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	welcomeMessage := "Welcome to the server!\nEnter 'exit' to close the connection.\nEnter your input: "
	conn.Write([]byte(welcomeMessage))

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		clientInput := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Received:", clientInput)

		if clientInput == "exit" {
			fmt.Println("Closing connection as per client request.")
			return
		}
	}
}
