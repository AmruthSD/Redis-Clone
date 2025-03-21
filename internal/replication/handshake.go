package replication

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
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

func UpdateOffset(parts []string) {
	tt := 0
	for _, v := range parts {
		tt += len(v)
	}
	Metadata.MasterReplOffset += tt
	if Metadata.Role == "master" {
		storage.InsertCommand(parts, tt)
	}
}

func MakeHandShake() (net.Conn, error) {
	conn, err := net.Dial("tcp", Metadata.MasterAddress)
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
	off, err := strconv.Atoi(parts[2])
	if err == nil {
		Metadata.MasterReplOffset = off
	} else {
		Metadata.MasterReplOffset = 0
	}

	file, err := os.OpenFile("../"+config.RedisConfig.DbFileName+"/"+config.RedisConfig.DbFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}

	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error receiving file:", err)
		return nil, err
	}
	fmt.Println("File received successfully!")

	file.Close()
	storage.Reader()

	conn.Write([]byte(fmt.Sprintf("PSYNC %s %d", Metadata.MasterReplid, Metadata.MasterReplOffset)))
	return conn, nil
}
