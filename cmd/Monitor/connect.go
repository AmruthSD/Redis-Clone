package main

import (
	"fmt"
	"net"
	"time"
)

var ConnctionMap map[string]string = make(map[string]string)
var ConnectionChan = make(map[string]chan string)
var ConnectionCount = make(map[string]int)

func handle_read(buf []byte, con net.Conn) error {
	bytesRead, err := con.Read(buf)
	if err != nil {
		return err
	}
	message := string(buf[:bytesRead])
	fmt.Println("Message received:", message)
	switch message {
	case "PING":
		con.Write([]byte("PONG"))
	case "SEND MASTER CONTACT":
		if MasterKey == con.RemoteAddr().String() {
			con.Write([]byte("YOU ARE THE MASTER"))
		} else {
			con.Write([]byte(ConnctionMap[MasterKey]))
		}
	case "SEND SLAVE CONTACT":
		go SendSlave(con)
	case "CLOSED A CONNECTION":
		ConnectionCount[con.RemoteAddr().String()]--
	default:
		ConnctionMap[con.RemoteAddr().String()] = message
		con.Write([]byte("REGISTERED"))
	}
	return nil
}

func Connect(conn net.Conn) {
	ConnctionMap[conn.RemoteAddr().String()] = "-1"
	ConnectionChan[conn.RemoteAddr().String()] = make(chan string)
	ConnectionCount[conn.RemoteAddr().String()] = 0
	defer conn.Close()
	defer delete(ConnctionMap, conn.RemoteAddr().String())
	defer delete(ConnectionChan, conn.RemoteAddr().String())
	defer close(ConnectionChan[conn.RemoteAddr().String()])
	defer delete(ConnectionCount, conn.RemoteAddr().String())
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
