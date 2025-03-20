package connection

import (
	"errors"
	"fmt"
	"net"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

const MonitorPort = "7000"

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
	fmt.Println("Received from Master:", message)

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
	fmt.Println("Received from Master:", message)

	return conn, nil
}

func HandleMonitorConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		bytesRead, err := conn.Read(buf)
		message := string(buf[:bytesRead])
		if err != nil {
			fmt.Println("Error reading response:", err)

		} else if message == "YOU ARE THE MASTER" {
			replication.Metadata.Role = "master"
			Master_Init()
		} else {
			replication.Metadata.Role = "slave"
			replication.Metadata.MasterAddress = message
			Slave_Init()
		}
		fmt.Println("Received from Master:", message)
	}
}
