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

	for {
		fmt.Print("Enter your input: ")
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
