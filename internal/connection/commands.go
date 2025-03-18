package connection

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func handle_ping(parts []string, conn net.Conn) {
	conn.Write([]byte("PONG\n"))
}

func handle_echo(parts []string, conn net.Conn) {
	if len(parts) < 2 {
		conn.Write([]byte("Argument Count Not Right"))
	} else {
		conn.Write([]byte(strings.Join(parts[1:], " ")))
	}
}

func handle_set(parts []string, conn net.Conn) {

	if len(parts) == 3 {
		storage.SetValue(parts[1], parts[2], -1)
	} else if len(parts) == 5 && parts[3] == "PX" {
		ext, _ := strconv.ParseInt(parts[4], 10, 64)
		ti := time.Now().UnixMilli() + ext
		storage.SetValue(parts[1], parts[2], ti)
	} else if len(parts) != 3 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	}
	conn.Write([]byte("Done"))
}

func handle_get(parts []string, conn net.Conn) {
	if len(parts) != 2 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	}
	val := storage.GetValue(parts[1])
	conn.Write([]byte(val))
}

func handle_config(parts []string, conn net.Conn) {
	if len(parts) != 3 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	} else if parts[1] == "GET" && parts[2] == "dir" {
		conn.Write([]byte(config.RedisConfig.Dir))
	} else if parts[1] == "GET" && parts[2] == "dbfilename" {
		conn.Write([]byte(config.RedisConfig.DbFileName))
	} else {
		conn.Write([]byte("Invalid Arguments"))
	}
}

func handle_keys(parts []string, conn net.Conn) {
	if len(parts) == 2 {
		s := storage.HasPrefix(parts[1])
		if len(s) == 0 {
			s = "NO KEYS FOUND"
		}
		conn.Write([]byte(s))
	} else {
		conn.Write([]byte("Argument Count Not Right"))
	}

}

func handle_info(parts []string, conn net.Conn) {
	if len(parts) == 2 && parts[1] == "replication" {
		s := "role:" + replication.Metadata.Role + "\n"
		s = s + "number of slaves:" + strconv.Itoa(replication.Metadata.NumberOfSlaves) + "\n"
		s = s + "master_replid:" + replication.Metadata.MasterReplid + "\n"
		s = s + "master_repl_offset:" + strconv.Itoa(replication.Metadata.MasterReplOffset) + "\n"
		conn.Write([]byte(s))
	} else {
		conn.Write([]byte("Argument Count Not Right"))
	}
}

func handle_replconf(parts []string, conn net.Conn) {
	if len(parts) == 3 && parts[1] == "listening-port" {
		num, _ := strconv.Atoi(parts[2])
		replication.SlavesConnections[conn.RemoteAddr().String()] = true
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
	} else {
		conn.Write([]byte("Error"))
	}
}
