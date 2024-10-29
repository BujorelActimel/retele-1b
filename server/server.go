package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Logger setup for different types of logs
var (
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
)

func initLoggers() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Open log file with current timestamp
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	logFile, err := os.OpenFile(fmt.Sprintf("logs/server_%s.log", currentTime),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Initialize different loggers
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	// Also write errors to file
	errorLog.SetOutput(logFile)
}

func main() {
	// Initialize loggers
	initLoggers()

	// Define command-line flags
	host := flag.String("host", "0.0.0.0", "Host to listen on")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", *host, *port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		errorLog.Fatal("Failed to start server:", err)
	}
	defer listen.Close()

	infoLog.Printf("Server started on %s", address)
	infoLog.Println("Available network interfaces:")
	printLocalIPs()

	// Start connection counter
	connectionCounter := 0

	for {
		conn, err := listen.Accept()
		if err != nil {
			errorLog.Println("Failed to accept connection:", err)
			continue
		}
		connectionCounter++
		go handleRequest(conn, connectionCounter)
	}
}

func printLocalIPs() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		errorLog.Println("Error getting local IPs:", err)
		return
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				infoLog.Printf("  - %s", ipnet.IP.String())
			}
		}
	}
}

func handleRequest(conn net.Conn, connID int) {
	// Record start time for request duration logging
	startTime := time.Now()

	// Get client information
	remoteAddr := conn.RemoteAddr().String()
	debugLog.Printf("[Conn #%d] New connection from %s", connID, remoteAddr)

	defer func() {
		conn.Close()
		duration := time.Since(startTime)
		debugLog.Printf("[Conn #%d] Connection closed. Duration: %v", connID, duration)
	}()

	reader := bufio.NewReader(conn)

	// Read problem number
	problemNum, err := reader.ReadString('\n')
	if err != nil {
		errorLog.Printf("[Conn #%d] Error reading problem number from %s: %v",
			connID, remoteAddr, err)
		return
	}
	problemNum = strings.TrimSpace(problemNum)

	infoLog.Printf("[Conn #%d] Client %s requested problem %s",
		connID, remoteAddr, problemNum)

	// Read input data and process based on problem number
	var result string

	switch problemNum {
	case "1":
		input, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading input: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 1 input: %s", connID, strings.TrimSpace(input))
		result = solveProblem1(input)

	case "2":
		input, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading input: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 2 input: %s", connID, strings.TrimSpace(input))
		result = solveProblem2(input)

	case "3":
		input, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading input: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 3 input: %s", connID, strings.TrimSpace(input))
		result = solveProblem3(input)

	case "4":
		str1, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading first string: %v", connID, err)
			return
		}
		str2, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading second string: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 4 inputs: %s, %s",
			connID, strings.TrimSpace(str1), strings.TrimSpace(str2))
		result = solveProblem4(str1, str2)

	case "5":
		input, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading input: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 5 input: %s", connID, strings.TrimSpace(input))
		result = solveProblem5(input)

	case "6":
		str, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading string: %v", connID, err)
			return
		}
		char, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading character: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 6 inputs: string=%s, char=%s",
			connID, strings.TrimSpace(str), strings.TrimSpace(char))
		result = solveProblem6(str, char)

	case "7":
		str, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading string: %v", connID, err)
			return
		}
		pos, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading position: %v", connID, err)
			return
		}
		length, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading length: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 7 inputs: string=%s, pos=%s, len=%s",
			connID, strings.TrimSpace(str), strings.TrimSpace(pos), strings.TrimSpace(length))
		result = solveProblem7(str, pos, length)

	case "8":
		arr1, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading first array: %v", connID, err)
			return
		}
		arr2, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading second array: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 8 inputs: arr1=%s, arr2=%s",
			connID, strings.TrimSpace(arr1), strings.TrimSpace(arr2))
		result = solveProblem8(arr1, arr2)

	case "9":
		arr1, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading first array: %v", connID, err)
			return
		}
		arr2, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading second array: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 9 inputs: arr1=%s, arr2=%s",
			connID, strings.TrimSpace(arr1), strings.TrimSpace(arr2))
		result = solveProblem9(arr1, arr2)

	case "10":
		str1, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading first string: %v", connID, err)
			return
		}
		str2, err := reader.ReadString('\n')
		if err != nil {
			errorLog.Printf("[Conn #%d] Error reading second string: %v", connID, err)
			return
		}
		debugLog.Printf("[Conn #%d] Problem 10 inputs: str1=%s, str2=%s",
			connID, strings.TrimSpace(str1), strings.TrimSpace(str2))
		result = solveProblem10(str1, str2)

	default:
		errorLog.Printf("[Conn #%d] Invalid problem number: %s", connID, problemNum)
		result = "Invalid problem number"
	}

	// Send result back to client
	debugLog.Printf("[Conn #%d] Sending result: %s", connID, result)
	conn.Write([]byte(result))
	infoLog.Printf("[Conn #%d] Request completed. Problem: %s, Duration: %v",
		connID, problemNum, time.Since(startTime))
}

// Sum of numbers
func solveProblem1(input string) string {
	numbers := strings.Fields(strings.TrimSpace(input))
	sum := 0
	for _, num := range numbers {
		n, err := strconv.Atoi(num)
		if err == nil {
			sum += n
		}
	}
	return fmt.Sprintf("Sum: %d", sum)
}

// Count spaces
func solveProblem2(input string) string {
	count := strings.Count(input, " ")
	return fmt.Sprintf("Number of spaces: %d", count)
}

// Reverse string
func solveProblem3(input string) string {
	runes := []rune(strings.TrimSpace(input))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Merge sorted strings
func solveProblem4(str1, str2 string) string {
	s1 := strings.Fields(strings.TrimSpace(str1))
	s2 := strings.Fields(strings.TrimSpace(str2))
	merged := make([]string, 0, len(s1)+len(s2))
	i, j := 0, 0

	for i < len(s1) && j < len(s2) {
		if s1[i] <= s2[j] {
			merged = append(merged, s1[i])
			i++
		} else {
			merged = append(merged, s2[j])
			j++
		}
	}

	merged = append(merged, s1[i:]...)
	merged = append(merged, s2[j:]...)
	return strings.Join(merged, " ")
}

// Divisors
func solveProblem5(input string) string {
	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "Invalid number"
	}

	var divisors []string
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			divisors = append(divisors, strconv.Itoa(i))
		}
	}
	return strings.Join(divisors, " ")
}

// Find character positions
func solveProblem6(str, char string) string {
	str = strings.TrimSpace(str)
	char = strings.TrimSpace(char)
	if len(char) == 0 {
		return "No character provided"
	}

	var positions []string
	for i, c := range str {
		if string(c) == char[:1] {
			positions = append(positions, strconv.Itoa(i))
		}
	}
	return strings.Join(positions, " ")
}

// Substring
func solveProblem7(str, posStr, lenStr string) string {
	str = strings.TrimSpace(str)
	pos, err1 := strconv.Atoi(strings.TrimSpace(posStr))
	length, err2 := strconv.Atoi(strings.TrimSpace(lenStr))

	if err1 != nil || err2 != nil {
		return "Invalid position or length"
	}

	if pos < 0 || pos >= len(str) || length < 0 {
		return "Invalid parameters"
	}

	end := pos + length
	if end > len(str) {
		end = len(str)
	}

	return str[pos:end]
}

// Common numbers
func solveProblem8(arr1, arr2 string) string {
	nums1 := strings.Fields(strings.TrimSpace(arr1))
	nums2 := strings.Fields(strings.TrimSpace(arr2))

	set := make(map[string]bool)
	for _, num := range nums1 {
		set[num] = true
	}

	var common []string
	for _, num := range nums2 {
		if set[num] {
			common = append(common, num)
			delete(set, num) // Avoid duplicates
		}
	}
	return strings.Join(common, " ")
}

// Numbers in first but not in second array
func solveProblem9(arr1, arr2 string) string {
	nums1 := strings.Fields(strings.TrimSpace(arr1))
	nums2 := strings.Fields(strings.TrimSpace(arr2))

	set := make(map[string]bool)
	for _, num := range nums2 {
		set[num] = true
	}

	var diff []string
	for _, num := range nums1 {
		if !set[num] {
			diff = append(diff, num)
		}
	}
	return strings.Join(diff, " ")
}

// Most common character at same positions
func solveProblem10(str1, str2 string) string {
	str1 = strings.TrimSpace(str1)
	str2 = strings.TrimSpace(str2)

	if len(str1) == 0 || len(str2) == 0 {
		return "Empty string(s)"
	}

	charCount := make(map[rune]int)
	minLen := len(str1)
	if len(str2) < minLen {
		minLen = len(str2)
	}

	for i := 0; i < minLen; i++ {
		if str1[i] == str2[i] {
			charCount[rune(str1[i])]++
		}
	}

	maxChar := ' '
	maxCount := 0
	for char, count := range charCount {
		if count > maxCount {
			maxCount = count
			maxChar = char
		}
	}

	if maxCount == 0 {
		return "No characters match"
	}

	return fmt.Sprintf("Character: %c, Count: %d", maxChar, maxCount)
}
