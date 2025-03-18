package connection

import (
	"fmt"
	"net"
	"strings"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

func HandleConnection(con net.Conn) {
	defer delete(replication.SlavesConnections, con.RemoteAddr().String())
	defer con.Close()

	buf := make([]byte, 1024)

	for {
		bytesRead, err := con.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading from connection:", err)
			} else if err.Error() == "EOF" {
				fmt.Println("Connection Closed")
			}
			return
		}

		message := string(buf[:bytesRead])
		fmt.Println("Message received:", message)

		parts := strings.Fields(message)

		handler, err := Parse(parts)
		if err == nil {
			handler(parts, con)
		} else {
			con.Write([]byte(err.Error()))
		}
	}
}
