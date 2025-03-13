package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/connection"
	"github.com/AmruthSD/Redis-Clone/internal/replication"
)

func main() {

	config.LoadConfig()

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.RedisConfig.Port))
	if err != nil {
		fmt.Println("Failed to bind to port ", config.RedisConfig.Port)
		os.Exit(1)
	}

	defer l.Close()

	if config.RedisConfig.ReplicaOf != "" {
		replication.Metadata.Role = "slave"
	}

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
