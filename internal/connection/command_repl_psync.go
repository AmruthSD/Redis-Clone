package connection

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func handle_replconf(parts []string, conn net.Conn) {
	if len(parts) == 3 && parts[1] == "listening-port" {
		num, _ := strconv.Atoi(parts[2])
		replication.SlavesConnections[conn.RemoteAddr().String()] = true
		replication.ConnectionChannels[conn.RemoteAddr().String()] = make(chan string, 10)
		if num <= 1<<16 {
			conn.Write([]byte("OK"))
		} else {
			conn.Write([]byte("Error"))
		}
	} else if len(parts) == 3 && parts[1] == "capa" && parts[2] == "psync2" {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte("Error"))
	}
}

func handle_psync(parts []string, conn net.Conn) {
	if len(parts) == 3 && parts[1] == "?" && parts[2] == "-1" {
		conn.Write([]byte(fmt.Sprintf("FULLRESYNC %s 0", replication.Metadata.MasterReplid)))
		/*
			file, err := os.Open(config.RedisConfig.Dir + "/" + config.RedisConfig.DbFileName)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()
			_, err = io.Copy(conn, file)
			if err != nil {
				fmt.Println("Error sending file:", err)
			}
			fmt.Println("File sent successfully!")
		*/
	} else if parts[1] == replication.Metadata.MasterReplid {
		slave_off, err := strconv.Atoi(parts[2])
		if err != nil {
			conn.Write([]byte("Error"))
			return
		}
		by := replication.Metadata.MasterReplOffset - slave_off
		buf, err := storage.LastFewCommands(by)
		if err != nil {
			conn.Write([]byte("Error"))
			return
		}
		conn.Write([]byte("OK"))
		for _, val := range buf {
			replication.ConnectionChannels[conn.RemoteAddr().String()] <- strings.Join(val, " ")
		}
	} else {
		conn.Write([]byte("Error"))
	}
}
