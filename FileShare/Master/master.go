package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting the server...")
	ln, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal("Error accepting connection:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to client")

	/************SEND OR RECIEVE***************/
	// Read user input (send or receive)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 'send' or 'receive': ")
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)

	if command == "send" {
		// Assuming you are receiving processed files back from slaves before combining them
		combineFiles("combinedfile.txt", "part1.txt", "part2.txt")
		sendFileToClient("combinedfile.txt", conn)
		fmt.Println("Sent combined file to client")
		conn.Close()
		fmt.Println("Connection closed")
	} else if command == "receive" {
		newFile, err := os.Create("newfile.txt")
		if err != nil {
			log.Fatal("Error creating new file:", err)
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, conn)
		if err != nil {
			log.Fatal("Error receiving file from client:", err)
		}
		fmt.Println("Received file from client")

		file, err := os.Open("newfile.txt")
		if err != nil {
			log.Fatal("Error opening new file:", err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal("Error getting file info:", err)
		}

		partSize := fileInfo.Size() / 2
		part1File, err := os.Create("part1.txt")
		if err != nil {
			log.Fatal("Error creating part1 file:", err)
		}
		defer part1File.Close()

		_, err = io.CopyN(part1File, file, partSize)
		if err != nil && err != io.EOF {
			log.Fatal("Error copying first part:", err)
		}

		part2File, err := os.Create("part2.txt")
		if err != nil {
			log.Fatal("Error creating part2 file:", err)
		}
		defer part2File.Close()

		_, err = io.Copy(part2File, file)
		if err != nil {
			log.Fatal("Error copying second part:", err)
		}
		fmt.Println("Divided file into two parts")
		// Assume both slaves are on the same machine, different ports
		sendFileToSlave("part1.txt", "slave1LocalHost:50001")
		sendFileToSlave("part2.txt", "slave2LocalHost:50002")
	}

	/**********************************/

	/************************************/

}

func sendFileToSlave(filename string, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("Error dialing slave:", err)
	}
	defer conn.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening part file:", err)
	}
	defer file.Close()

	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Error sending file to slave:", err)
	}
	fmt.Printf("File %s sent to slave at %s\n", filename, address)
}

func combineFiles(combinedFilename string, part1Filename string, part2Filename string) {
	combinedFile, err := os.Create(combinedFilename)
	if err != nil {
		log.Fatal("Error creating combined file:", err)
	}
	defer combinedFile.Close()

	part1File, err := os.Open(part1Filename)
	if err != nil {
		log.Fatal("Error opening part1 file:", err)
	}
	defer part1File.Close()

	part2File, err := os.Open(part2Filename)
	if err != nil {
		log.Fatal("Error opening part2 file:", err)
	}
	defer part2File.Close()

	_, err = io.Copy(combinedFile, part1File)
	if err != nil {
		log.Fatal("Error combining part1:", err)
	}

	_, err = io.Copy(combinedFile, part2File)
	if err != nil {
		log.Fatal("Error combining part2:", err)
	}
	fmt.Println("Combined both parts into a single file")
}

func sendFileToClient(filename string, conn net.Conn) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening combined file:", err)
	}
	defer file.Close()

	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Error sending combined file to client:", err)
	}
}
