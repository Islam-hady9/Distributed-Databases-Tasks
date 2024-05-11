// Slave.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Listen on a port
	ln, _ := net.Listen("tcp", ":50001")
	defer ln.Close()

	fmt.Println("Slave is waiting for connection from master...")

	// Accept connection from master
	conn, _ := ln.Accept()
	defer conn.Close()
	fmt.Println("Connected to master")

	// Receive part from master
	partFile, err := os.Create("part.txt")
	if err != nil {
		fmt.Println("Error creating part file:", err)
		return
	}
	defer partFile.Close()

	// Copy the received data to the part file
	n, err := io.Copy(partFile, conn)
	if err != nil {
		fmt.Println("Error copying data:", err)
		return
	}
	fmt.Printf("Received %d bytes from master\n", n) // Print the number of bytes received

	// Send acknowledgment to master
	conn.Write([]byte("ACK")) // Send a simple acknowledgment

	// Wait for master to close the connection
	fmt.Println("Waiting for master to close the connection...")
	conn.Close()
	fmt.Println("Master connection closed")
}
