package main

import (
	"fmt"
	"net"
	"os"
)

func handleConnection(con net.Conn) {
	defer con.Close()

	buf := make([]byte, 1024)

	for {
		bytesRead, err := con.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading from connection:", err)
			} else if err.Error() == "EOF" {
				fmt.Println("Connection Closed")
			}
			return
		}

		message := string(buf[:bytesRead])
		fmt.Println("Message received:", message)
		if message == "PING\r\n" {
			con.Write([]byte("+PONG\r\n"))
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(con)
	}
}
