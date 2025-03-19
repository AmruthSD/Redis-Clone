package replication

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/config"
)

const alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func new_replication_id() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 40; i++ {
		s = s + string(alphanumeric[rand.Int()%len(alphanumeric)])
	}
	return s
}

func MakeHandShake() (net.Conn, error) {
	conn, err := net.Dial("tcp", "0.0.0.0:"+config.RedisConfig.ReplicaOf)
	if err != nil {
		fmt.Println("Error Connecting to the Master:", err)
		return nil, err
	}

	conn.Write([]byte("PING"))
	buf := make([]byte, 1024)
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

	conn.Write([]byte(fmt.Sprintf("REPLCONF listening-port %d", config.RedisConfig.Port)))
	bytesRead, err = conn.Read(buf)
	message = string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "OK" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Master:", message)

	conn.Write([]byte("REPLCONF capa psync2"))
	bytesRead, err = conn.Read(buf)
	message = string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "OK" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Master:", message)

	conn.Write([]byte("PSYNC ? -1"))
	bytesRead, err = conn.Read(buf)
	message = string(buf[:bytesRead])
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}
	fmt.Println("Received from Master:", message)
	parts := strings.Split(message, " ")
	fmt.Println("Handshake Done")
	Metadata.MasterReplid = parts[1]
	Metadata.MasterReplOffset = 0
	/*
		file, err := os.Create("received.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()

		bytesRead, err = io.Copy(file, conn)
		if err != nil {
			fmt.Println("Error receiving file:", err)
			return err
		}
		fmt.Println("File received successfully!")
	*/

	return conn, nil
}
