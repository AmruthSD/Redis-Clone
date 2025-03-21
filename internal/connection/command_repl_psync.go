package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func handle_replconf(parts []string, conn net.Conn) {
	if len(parts) == 3 && parts[1] == "listening-port" {
		num, _ := strconv.Atoi(parts[2])
		replication.SlavesConnections[conn.RemoteAddr().String()] = true
		replication.ConnectionChannels[conn.RemoteAddr().String()] = make(chan string, 10)
		if num <= 1<<16 {
			conn.Write([]byte("OK\n"))
		} else {
			conn.Write([]byte("Error\n"))
		}
	} else if len(parts) == 3 && parts[1] == "capa" && parts[2] == "psync2" {
		conn.Write([]byte("OK\n"))
	} else {
		conn.Write([]byte("Error\n"))
	}
}

func handle_psync(parts []string, conn net.Conn) {
	if len(parts) == 3 && parts[1] == "?" && parts[2] == "-1" {
		conn.Write([]byte(fmt.Sprintf("FULLRESYNC %s %d\n", replication.Metadata.MasterReplid, storage.SlaveOffsetVal)))

		file, err := os.OpenFile("./"+config.RedisConfig.Dir+"/"+config.RedisConfig.DbFileName, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Error reading from file:", err)
				return
			}

			if _, err := conn.Write(buf[:n]); err != nil {
				fmt.Println("Error sending file:", err)
				return
			}
		}
		conn.Write([]byte("FILE SENT SUCCESFULLY\n"))
		fmt.Println("File sent successfully!")

	} else if parts[1] == replication.Metadata.MasterReplid {
		slave_off, err := strconv.Atoi(parts[2])
		if err != nil {
			conn.Write([]byte("Error\n"))
			return
		}
		by := storage.MasterReplOffset - slave_off
		buf, err := storage.LastFewCommands(by)
		if err != nil {
			conn.Write([]byte("Error\n"))
			return
		}
		conn.Write([]byte("OK\n"))
		for _, val := range buf {
			replication.ConnectionChannels[conn.RemoteAddr().String()] <- strings.Join(val, " ")
		}
	} else {
		conn.Write([]byte("Error\n"))
	}
}
