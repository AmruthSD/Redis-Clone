package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	serverAddr := "0.0.0.0:6379"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to", serverAddr)

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		num = 1
	}
	for range num {
		msg := "PING\r\n"
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println("Received:", string(buf))
		time.Sleep(1 * time.Second)
	}
}
