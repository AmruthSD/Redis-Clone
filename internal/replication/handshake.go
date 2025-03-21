package replication

import (
	"bufio"
	"errors"
	"fmt"
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
	storage.MasterReplOffset += tt
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

	conn.Write([]byte("PING\n"))
	reader := bufio.NewReader(conn)

	message, err := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "PONG" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Master:", message)

	conn.Write([]byte(fmt.Sprintf("REPLCONF listening-port %d\n", config.RedisConfig.Port)))
	message, err = reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "OK" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Master:", message)

	conn.Write([]byte("REPLCONF capa psync2\n"))
	message, err = reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	} else if message != "OK" {
		fmt.Println("Incorrect Response:", message)
		return nil, errors.New("incorrect response")
	}
	fmt.Println("Received from Master:", message)

	conn.Write([]byte("PSYNC ? -1\n"))
	message, err = reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}
	fmt.Println("Received from Master:", message)
	parts := strings.Split(message, " ")

	Metadata.MasterReplid = parts[1]
	off, err := strconv.Atoi(parts[2])
	if err == nil {
		storage.MasterReplOffset = off
	} else {
		storage.MasterReplOffset = 0
	}
	fmt.Println("Starting Reading File")
	file, err := os.OpenFile("./"+config.RedisConfig.Dir+"/"+config.RedisConfig.DbFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}
	fmt.Println("File opened successfully!")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if scanner.Text() == "FILE SENT SUCCESFULLY" {
			break
		}
		_, err = file.WriteString(scanner.Text() + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
		return nil, err
	}
	fmt.Println("File received successfully!")

	file.Close()
	storage.Reader()

	conn.Write([]byte(fmt.Sprintf("PSYNC %s %d\n", Metadata.MasterReplid, storage.MasterReplOffset)))
	fmt.Println("Handshake Done")
	return conn, nil
}
