// Client.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the master
	conn, _ := net.Dial("tcp", "masterLocalHost:50000")
	defer conn.Close()

	// Read user input (send or receive)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 'send' or 'receive': ")
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)

	if command == "send" {
		// Send file
		file, _ := os.Open("file.txt")
		io.Copy(conn, file)
		file.Close()
		fmt.Println("Sent file to master")
	} else if command == "receive" {
		// Receive the combined file from the master
		combinedFile, err := os.Create("received_combinedfile.txt")
		if err != nil {
			fmt.Println("Error creating received combined file:", err)
			return
		}
		defer combinedFile.Close()

		io.Copy(combinedFile, conn)
		fmt.Println("Received combined file from master")
	}

}
