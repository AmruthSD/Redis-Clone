package connection

import (
	"fmt"
	"net"
	"strings"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

var ConnectionChannels map[string]chan int

func handle_accepet_request(buf []byte, con net.Conn) error {
	bytesRead, err := con.Read(buf)
	if err != nil {
		return err
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

	return nil
}

func HandleConnection(con net.Conn) {
	defer con.Close()
	defer delete(replication.SlavesConnections, con.RemoteAddr().String())
	defer delete(ConnectionChannels, con.RemoteAddr().String())

	buf := make([]byte, 1024)

	for {

		if replication.SlavesConnections[con.RemoteAddr().String()] {
			select {
			case val := <-ConnectionChannels[con.RemoteAddr().String()]:
				if val == 1 {

				} else {

				}
			default:
				err := handle_accepet_request(buf, con)
				if err != nil {
					if err.Error() != "EOF" {
						fmt.Println("Error reading from connection:", err)
					} else if err.Error() == "EOF" {
						fmt.Println("Connection Closed")
					}
					return
				}
			}
		} else {
			err := handle_accepet_request(buf, con)
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println("Error reading from connection:", err)
				} else if err.Error() == "EOF" {
					fmt.Println("Connection Closed")
				}
				return
			}
		}

	}
}
