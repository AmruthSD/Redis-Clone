package main

import (
	"fmt"
	"net"
	"os"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/connection"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	config.LoadConfig()
	// storage.ReadFile(config.RedisConfig.Dir + "/" + config.RedisConfig.DbFileName)

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go connection.HandleConnection(con)
	}
}
