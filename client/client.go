package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run client.go <host> <port> <problem_number>")
		os.Exit(1)
	}

	host := os.Args[1]
	port := os.Args[2]
	problemNum := os.Args[3]

	tcpServer, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Send problem number first
	_, err = conn.Write([]byte(problemNum + "\n"))
	if err != nil {
		fmt.Println("Write problem number failed:", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// Handle different input requirements based on problem number
	switch problemNum {
	case "4", "8", "9", "10":
		fmt.Println("Enter first string/sequence:")
		str1, _ := reader.ReadString('\n')
		fmt.Println("Enter second string/sequence:")
		str2, _ := reader.ReadString('\n')
		message := str1 + str2
		_, err = conn.Write([]byte(message))
	case "6":
		fmt.Println("Enter string:")
		str, _ := reader.ReadString('\n')
		fmt.Println("Enter character to search:")
		char, _ := reader.ReadString('\n')
		message := str + char
		_, err = conn.Write([]byte(message))
	case "7":
		fmt.Println("Enter string:")
		str, _ := reader.ReadString('\n')
		fmt.Println("Enter start position:")
		pos, _ := reader.ReadString('\n')
		fmt.Println("Enter length:")
		length, _ := reader.ReadString('\n')
		message := strings.Join([]string{str, pos, length}, "\n")
		_, err = conn.Write([]byte(message))
	default:
		fmt.Println("Enter your input:")
		message, _ := reader.ReadString('\n')
		_, err = conn.Write([]byte(message))
	}

	if err != nil {
		fmt.Println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// Buffer to get data
	received := make([]byte, 1024)
	n, err := conn.Read(received)
	if err != nil {
		fmt.Println("Read data failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Server response:", string(received[:n]))
}
