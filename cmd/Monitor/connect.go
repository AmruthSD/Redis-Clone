package main

import (
	"fmt"
	"net"
	"time"
)

var ConnctionMap map[string]bool = make(map[string]bool)
var ConnectionChan = make(map[string]chan string)

func handle_read(buf []byte, con net.Conn) error {
	bytesRead, err := con.Read(buf)
	if err != nil {
		return err
	}
	message := string(buf[:bytesRead])
	fmt.Println("Message received:", message)

	return nil
}

func Connect(conn net.Conn) {
	ConnctionMap[conn.RemoteAddr().String()] = true
	ConnectionChan[conn.RemoteAddr().String()] = make(chan string)
	defer conn.Close()
	defer delete(ConnctionMap, conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		select {
		case val := <-ConnectionChan[conn.RemoteAddr().String()]:
			conn.Write([]byte(val))
		default:
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			err := handle_read(buf, conn)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				} else if err.Error() != "EOF" {
					fmt.Println("Error reading from connection:", err)
				} else {
					fmt.Println("Connection Closed")
				}
				return
			}
		}
	}
}
