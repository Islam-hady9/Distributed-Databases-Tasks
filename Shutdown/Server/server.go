package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	slaveConnections    = make(map[string]net.Conn)
	slaveConnectionsMux sync.Mutex
)

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "serverlocalhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Master device is listening for connections...")

	// Start a goroutine to handle shutdown signal
	go handleShutdownSignal()

	// Accept connections from slaves
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Slave connected:", conn.RemoteAddr())

		slaveConnectionsMux.Lock()
		slaveConnections[conn.RemoteAddr().String()] = conn
		slaveConnectionsMux.Unlock()

		go handleSlaveConnection(conn) // Handle slave connection in a separate goroutine
	}
}

func handleSlaveConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Handling slave connection:", conn.RemoteAddr())

	// Read and write data with the slave device
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from slave:", err)
			break
		}
		log.Println("Received data from slave:", string(buffer[:n]))

		// Check if the received message is the shutdown signal from the master
		if string(buffer[:n]) == "shutdown" {
			log.Println("Received shutdown signal from master. Shutting down slave:", conn.RemoteAddr())
			break
		}

		// Example: Write to the slave
		message := []byte("Hello from the master!")
		_, err = conn.Write(message)
		if err != nil {
			log.Println("Error writing to slave:", err)
			break
		}
		log.Println("Sent data to slave")
	}

	// Remove slave connection from the map
	slaveConnectionsMux.Lock()
	delete(slaveConnections, conn.RemoteAddr().String())
	slaveConnectionsMux.Unlock()
}

func handleShutdownSignal() {
	for {
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		if input == "shutdown" {
			log.Println("Shutting down all slave connections...")

			slaveConnectionsMux.Lock()
			for _, conn := range slaveConnections {
				_, err := conn.Write([]byte("shutdown"))
				if err != nil {
					log.Println("Error sending shutdown signal to slave:", err)
				}
				conn.Close()
			}
			slaveConnections = make(map[string]net.Conn)
			slaveConnectionsMux.Unlock()

			log.Println("All slave connections have been closed.")
		}
	}
}
