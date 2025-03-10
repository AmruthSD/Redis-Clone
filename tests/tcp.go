package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "0.0.0.0:6379"
	con, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer con.Close()
	fmt.Println("Connected to", serverAddr)
}
