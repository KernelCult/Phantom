package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type ListenerInfo struct {
	User       string
	ListenerID int
}

var (
	activeUsers      sync.Map
	activeListeners  sync.Map
	listenerCounter  int
	userListenerInfo sync.Map
)

func main() {
	listener, err := net.Listen("tcp", "192.168.56.138:8080")
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

		if strings.HasPrefix(clientInput, "list users") {
			listActiveUsers(conn)
		} else if strings.HasPrefix(clientInput, "start listener") {
			startListener(clientInput, username)
			conn.Write([]byte("\n[server] Listener started successfully.\n" + username + "> "))
		} else if strings.HasPrefix(clientInput, "list listener") {
			listActiveListeners(conn)
		} else {
			conn.Write([]byte("\n[server] Invalid command.\n" + username + "> "))
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

func startListener(command, username string) {
	// Extract port from the command
	parts := strings.Split(command, " ")
	if len(parts) != 3 {
		fmt.Println("Invalid 'start listener' command format.")
		return
	}

	port := parts[2]

	// Start a listener on the specified port
	go func(port, username string) {
		listener, err := net.Listen("tcp", "192.168.56.138:"+port)
		if err != nil {
			fmt.Printf("Error starting listener on port %s: %s\n", port, err)
			return
		}
		defer listener.Close()

		id := generateListenerID()
		fmt.Printf("\nListener %d started on port %s by user '%s'\n", id, port, username)

		activeListeners.Store(id, listener)
		userListenerInfo.Store(id, ListenerInfo{User: username, ListenerID: id})
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection on listener %d: %s\n", id, err)
				return
			}
			fmt.Printf("Listener %d accepted connection.\n", id)
			go handleListenerConnection(id, conn)
		}
	}(port, username)
}

func handleListenerConnection(listenerID int, conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("\nListener %d connection closed.\n", listenerID)
			activeListeners.Delete(listenerID)
			userListenerInfo.Delete(listenerID)
			return
		}

		clientInput := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("\nListener %d received: %s\n", listenerID, clientInput)
	}
}

func listActiveListeners(conn net.Conn) {
	conn.Write([]byte("Active Listeners:\n"))
	activeListeners.Range(func(key, value interface{}) bool {
		listenerID := key.(int)
		userInfo, ok := userListenerInfo.Load(listenerID)
		if ok {
			conn.Write([]byte(fmt.Sprintf("- Listener %d started by user '%s'\n", listenerID, userInfo.(ListenerInfo).User)))
		} else {
			conn.Write([]byte(fmt.Sprintf("- Listener %d (No user information available)\n", listenerID)))
		}
		return true
	})
}

func getUsernameForListener(listenerID int) string {
	userInfo, ok := userListenerInfo.Load(listenerID)
	if ok {
		return userInfo.(ListenerInfo).User
	}
	return "UnknownUser"
}

func generateListenerID() int {
	listenerCounter++
	return listenerCounter
}
