package connection

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

var MyReceivingAddress string

func handle_accepet_request(con net.Conn) error {
	reader := bufio.NewReader(con)

	message, err := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		return err
	}

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

	for {

		if replication.SlavesConnections[con.RemoteAddr().String()] {

			select {
			case msg := <-replication.ConnectionChannels[con.RemoteAddr().String()]:
				fmt.Println("Sending to Slave:", msg)
				con.Write([]byte(msg + "\n"))
			default:
				con.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
				err := handle_accepet_request(con)
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
			err := handle_accepet_request(con)
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

	for {

		err := handle_accepet_request(con)
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
