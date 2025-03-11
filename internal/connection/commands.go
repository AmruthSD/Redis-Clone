package connection

import (
	"net"
	"strings"

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
	if len(parts) != 3 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	}
	storage.SetValue(parts[1], parts[2])
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
