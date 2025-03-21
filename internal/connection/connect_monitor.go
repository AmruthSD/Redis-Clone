package connection

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

const MonitorPort = "7000"

var Monitor_Conn net.Conn

func Connect_Monitor() (net.Conn, error) {
	conn, err := net.Dial("tcp", "0.0.0.0:"+MonitorPort)
	if err != nil {
		fmt.Println("Error Connecting to the Monitor:", err)
		return nil, err
	}

	buf := make([]byte, 1024)

	conn.Write([]byte("PING"))
	bytesRead, err := conn.Read(buf)
	message := string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "PONG" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Monitor:", message)

	conn.Write([]byte(MyReceivingAddress))
	bytesRead, err = conn.Read(buf)
	message = string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "REGISTERED" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Monitor:", message)

	conn.Write([]byte("SEND MASTER CONTACT"))
	bytesRead, err = conn.Read(buf)
	message = string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message == "YOU ARE THE MASTER" {
		replication.Metadata.Role = "master"
		Master_Init()
	} else {
		replication.Metadata.Role = "slave"
		replication.Metadata.MasterAddress = message
		Slave_Init()
	}
	fmt.Println("Received from Monitor:", message)

	return conn, nil
}

func HandleMonitorConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		bytesRead, err := conn.Read(buf)
		message := string(buf[:bytesRead])
		if err != nil && err.Error() == "EOF" {
			fmt.Println("MONITOR CLOSED")
			os.Exit(3)
		} else if err != nil {
			fmt.Println("Error reading response:", err)

		} else if message == "YOU ARE THE MASTER" {
			replication.Metadata.Role = "master"
			Master_Init()
		} else {
			replication.Metadata.Role = "slave"
			replication.Metadata.MasterAddress = message
			Slave_Init()
		}
		fmt.Println("Received from Monitor:", message)
	}
}
