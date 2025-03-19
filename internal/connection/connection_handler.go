package connection

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

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
	defer delete(replication.ConnectionChannels, con.RemoteAddr().String())

	buf := make([]byte, 1024)

	for {

		if replication.SlavesConnections[con.RemoteAddr().String()] {

			select {
			case msg := <-replication.ConnectionChannels[con.RemoteAddr().String()]:
				fmt.Println("Sending to Slave:", msg)
				con.Write([]byte(msg))
			default:
				con.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
				err := handle_accepet_request(buf, con)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					} else if err.Error() != "EOF" {
						fmt.Println("Error reading from connection:", err)
					} else {
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

func HandleMasterConnection(con net.Conn) {
	defer con.Close()
	defer delete(replication.SlavesConnections, con.RemoteAddr().String())
	defer delete(replication.ConnectionChannels, con.RemoteAddr().String())

	buf := make([]byte, 1024)

	for {

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
