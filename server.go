package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var activeUsers sync.Map

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

	// welcomeMessage := "Welcome to the server!"
	// conn.Write([]byte(welcomeMessage))

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	username := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("User '%s' connected.\n", username)

	activeUsers.Store(username, struct{}{})

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("User '%s' disconnected.\n", username)
			activeUsers.Delete(username)
			return
		}

		clientInput := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("User '%s' sent: %s\n", username, clientInput)

		if clientInput == "list active_users" {
			listActiveUsers(conn)
		}
	}
}

func listActiveUsers(conn net.Conn) {
	conn.Write([]byte("Active Users:\n"))
	activeUsers.Range(func(key, value interface{}) bool {
		conn.Write([]byte(fmt.Sprintf("- %s\n", key)))
		return true
	})
}
