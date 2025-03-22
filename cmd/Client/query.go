package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func Query() {
	fmt.Printf("Query> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")

	if parts[0] == "EXIT" {
		os.Exit(1)
	}

	switch parts[0] {
	case "SET", "DEL":
		master_conn.Write([]byte(input + "\n"))
		ReadReply(master_conn)
	case "GET", "PING", "ECHO", "KEYS":
		slave_conn.Write([]byte(input + "\n"))
		ReadReply(slave_conn)
	default:
		master_conn.Write([]byte(input + "\n"))
		ReadReply(master_conn)
	}
}

func ReadReply(conn net.Conn) {
	buf := make([]byte, 1024)
	bytesRead, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			handshake(monitor_conn)
		} else {
			fmt.Println("Error:", err)
		}
		return
	}
	message := string(buf[:bytesRead])
	fmt.Println(message)
}
