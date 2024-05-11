// Master.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting the server ...")
	// Listen on a port
	ln, _ := net.Listen("tcp", ":50000")
	defer ln.Close()

	// Accept connection from client
	conn, _ := ln.Accept()
	defer conn.Close()
	fmt.Println("Connected to client")

	// Receive file from client
	newFile, _ := os.Create("newfile.txt")
	io.Copy(newFile, conn)
	newFile.Close()
	fmt.Println("Received file from client")

	// Divide the file into two parts
	file, _ := os.Open("newfile.txt")
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	partSize := fileSize / 2

	// Create two temporary files for the parts
	part1File, _ := os.Create("part1.txt")
	part2File, _ := os.Create("part2.txt")

	// Copy the first part
	io.CopyN(part1File, file, partSize)
	part1File.Close()

	// Copy the second part
	io.Copy(part2File, file)
	part2File.Close()
	fmt.Println("Divided file into two parts")

	// Distribute parts to slaves (simplified example)
	slave1Conn, _ := net.Dial("tcp", "localhost:50001")
	slave2Conn, _ := net.Dial("tcp", "localhost:50002")

	// Send part1 to slave1
	part1File.Seek(0, 0)
	io.Copy(slave1Conn, part1File)
	slave1Conn.Close()

	// Send part2 to slave2
	part2File.Seek(0, 0)
	io.Copy(slave2Conn, part2File)
	slave2Conn.Close()
	

	// Combine both parts into a single file
	combinedFile, err := os.Create("combinedfile.txt")
	if err != nil {
		fmt.Println("Error creating combined file:", err)
		return
	}
	defer combinedFile.Close()

	// Read part1 and part2
	part1File, err = os.Open("part1.txt")
	if err != nil {
		fmt.Println("Error opening part1 file:", err)
		return
	}
	defer part1File.Close()

	part2File, err = os.Open("part2.txt")
	if err != nil {
		fmt.Println("Error opening part2 file:", err)
		return
	}
	defer part2File.Close()

	// Copy both parts to the combined file
	io.Copy(combinedFile, part1File)
	io.Copy(combinedFile, part2File)
	fmt.Println("Combined both parts into a single file")

	// Print the size of the combined file
	combinedFile.Seek(0, 0)
	combinedFileSize, _ := combinedFile.Stat()
	fmt.Printf("Size of combined file: %d bytes\n", combinedFileSize.Size())

	// Send the combined file to the client
	combinedFile.Seek(0, 0)
	io.Copy(conn, combinedFile)
	fmt.Println("Sent combined file to client")

	// Close the connection after sending the file
	conn.Close()
	fmt.Println("Connection closed")

}