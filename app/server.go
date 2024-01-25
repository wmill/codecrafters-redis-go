package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.ToUpper(scanner.Text()) == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if strings.ToUpper(scanner.Text()) == "ECHO" {
			alpha, _ := regexp.Compile("^[a-zA-Z]")
			for scanner.Scan() {
				if alpha.MatchString(scanner.Text()) {
					conn.Write([]byte("+" + scanner.Text() + "\r\n"))
					break
				} else {
					fmt.Println(scanner.Text())
				}
			}
		}
	}
}