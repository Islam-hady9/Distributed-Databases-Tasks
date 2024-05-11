package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func readServerMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	data := make([]byte, 4096)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println(err)
			return
		}

		message := string(data[:n])
		log.Println("Received from server:", message)

		// Check for a special command, e.g., "shutdown"
		handleServerCommand(message)
	}
}

func sendUserMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		conn.Write([]byte(message))
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func handleServerCommand(command string) {
	command = strings.TrimSpace(command)

	if command == "shutdown" {
		log.Println("Received shutdown command from server. Shutting down the client...")
		shutdownClient()
	}
}

func shutdownClient() {
	var cmd *exec.Cmd

	if isWindows() {
		cmd = exec.Command("shutdown", "/s", "/t", "0")
	} else if isLinux() {
		cmd = exec.Command("shutdown", "-h", "now")
	} else {
		log.Println("Unsupported operating system for shutdown.")
		return
	}

	err := cmd.Run()
	if err != nil {
		log.Println("Error shutting down client:", err)
	}
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func isLinux() bool {
	return os.PathSeparator == '/' && os.PathListSeparator == ':'
}

func main() {
	conn, err := net.Dial("tcp", "192.168.15.120:8000")
	if nil != err {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go readServerMessages(conn, &wg)
	go sendUserMessages(conn, &wg)

	wg.Wait()
}
