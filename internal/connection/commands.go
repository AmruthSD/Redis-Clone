package connection

import (
	"net"
	"strconv"
	"strings"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func handle_ping(parts []string, conn net.Conn) {
	conn.Write([]byte("PONG\n"))
}

func handle_ok(parts []string, conn net.Conn) {

}

func handle_unknown(parts []string, conn net.Conn) {

}

func handle_echo(parts []string, conn net.Conn) {
	if len(parts) < 2 {
		conn.Write([]byte("Argument Count Not Right\n"))
	} else {
		conn.Write([]byte(strings.Join(parts[1:], " \n")))
	}
}

func handle_config(parts []string, conn net.Conn) {
	if len(parts) != 3 {
		conn.Write([]byte("Argument Count Not Right\n"))
		return
	} else if parts[1] == "GET" && parts[2] == "dir" {
		conn.Write([]byte(config.RedisConfig.Dir + "\n"))
	} else if parts[1] == "GET" && parts[2] == "dbfilename" {
		conn.Write([]byte(config.RedisConfig.DbFileName + "\n"))
	} else {
		conn.Write([]byte("Invalid Arguments\n"))
	}
}

func handle_keys(parts []string, conn net.Conn) {
	if len(parts) == 2 {
		s := storage.HasPrefix(parts[1])
		if len(s) == 0 {
			s = "NO KEYS FOUND"
		}
		conn.Write([]byte(s + "\n"))
	} else {
		conn.Write([]byte("Argument Count Not Right\n"))
	}

}

func handle_info(parts []string, conn net.Conn) {
	if len(parts) == 2 && parts[1] == "replication" {
		s := "role:" + replication.Metadata.Role + "\n"
		s = s + "number of slaves:" + strconv.Itoa(replication.Metadata.NumberOfSlaves) + "\n"
		s = s + "master_replid:" + replication.Metadata.MasterReplid + "\n"
		s = s + "master_repl_offset:" + strconv.Itoa(storage.MasterReplOffset) + "\n"
		conn.Write([]byte(s + "\n"))
	} else {
		conn.Write([]byte("Argument Count Not Right\n"))
	}
}
