package main

import (
	"fmt"
	"net"
	"os"
)

const MONITOR_PORT = "7000"

var master_conn net.Conn
var slave_conn net.Conn
var monitor_conn net.Conn

func handshake(monitor_conn net.Conn) error {
	buf := make([]byte, 1024)
	monitor_conn.Write([]byte("SEND MASTER CONTACT"))
	bytesRead, err := monitor_conn.Read(buf)
	if err != nil {
		return err
	}
	message := string(buf[:bytesRead])
	master_conn, err = net.Dial("tcp", message)
	if err != nil {
		fmt.Println("Error Connecting to the Master:", err)
		return err
	}

	monitor_conn.Write([]byte("SEND SLAVE CONTACT"))
	bytesRead, err = monitor_conn.Read(buf)
	if err != nil {
		return err
	}
	message = string(buf[:bytesRead])
	slave_conn, err = net.Dial("tcp", message)
	if err != nil {
		fmt.Println("Error Connecting to the Slave:", err)
		return err
	}
	return nil
}

func main() {
	monitor_conn, err := net.Dial("tcp", "0.0.0.0:"+MONITOR_PORT)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer monitor_conn.Close()
	fmt.Println("Connected to Monitor")
	err = handshake(monitor_conn)
	if err != nil {
		fmt.Println("Error handshake:", err)
		os.Exit(1)
	}
	for {
		Query()
	}
}
