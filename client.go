package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// Read and send the username to the server during the initial connection
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	conn.Write([]byte(username))
	fmt.Print("Connected to Nexus\n")
	go readServerResponse(conn) // Start a goroutine to continuously read from the server

	for {
		fmt.Print(username + "> ")
		userInput, _ := reader.ReadString('\n')

		// Remove newline character from the input
		userInput = strings.TrimSpace(userInput)

		// Send the input to the server
		_, err := conn.Write([]byte(userInput))
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		if userInput == "exit" {
			fmt.Println("Exiting client as per user request.")
			return
		}
	}
}

func readServerResponse(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		serverResponse := string(buffer[:n])
		fmt.Println(serverResponse)
	}
}
