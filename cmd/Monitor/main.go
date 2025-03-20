package main

import (
	"fmt"
	"net"
	"os"
)

var MasterKey string = ""

func main() {
	fmt.Println("Starting Monitor")
	l, err := net.Listen("tcp", "0.0.0.0:"+"7000")
	if err != nil {
		os.Exit(1)
	}
	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		if MasterKey == "" {
			MasterKey = con.RemoteAddr().String()
			go Connect(con)
		} else {
			go Connect(con)
		}
	}
}
