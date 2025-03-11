package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
	sl := append([]string{"GET"}, os.Args[1:]...)
	msg := strings.Join(sl, " ")
	_, err = conn.Write([]byte(msg))
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

}
