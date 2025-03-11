package connection

import (
	"net"
	"strings"
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
