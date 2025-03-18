package replication

import (
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

func MakeHandShake() error {
	conn, err := net.Dial("tcp", "0.0.0.0:"+config.RedisConfig.ReplicaOf)
	if err != nil {
		fmt.Println("Error Connecting to the Master:", err)
		return err
	}

	conn.Write([]byte("PING"))
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil || string(buf) != "PONG" {
		fmt.Println("Error reading response:", err)
		return err
	}
	fmt.Println("Received from Master:", string(buf))

	conn.Write([]byte(fmt.Sprintf("REPLCONF listening-port %d", config.RedisConfig.Port)))
	_, err = conn.Read(buf)
	if err != nil || string(buf) != "OK" {
		fmt.Println("Error reading response:", err)
		return err
	}
	fmt.Println("Received from Master:", string(buf))

	conn.Write([]byte("REPLCONF capa psync2"))
	_, err = conn.Read(buf)
	if err != nil || string(buf) != "OK" {
		fmt.Println("Error reading response:", err)
		return err
	}
	fmt.Println("Received from Master:", string(buf))

	conn.Write([]byte("PSYNC ? -1"))
	_, err = conn.Read(buf)
	if err != nil || string(buf) != "OK" {
		fmt.Println("Error reading response:", err)
		return err
	}
	fmt.Println("Received from Master:", string(buf))
	parts := strings.Split(string(buf), " ")

	Metadata.MasterReplid = parts[1]
	Metadata.MasterReplOffset = 0
	/*
		file, err := os.Create("received.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, conn)
		if err != nil {
			fmt.Println("Error receiving file:", err)
			return err
		}
		fmt.Println("File received successfully!")
	*/
	return nil
}
