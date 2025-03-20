package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/AmruthSD/Redis-Clone/internal/config"
	"github.com/AmruthSD/Redis-Clone/internal/connection"
	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func main() {

	config.LoadConfig()

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.RedisConfig.Port))
	if err != nil {
		fmt.Println("Failed to bind to port ", config.RedisConfig.Port)
		os.Exit(1)
	}

	defer l.Close()

	go storage.Single_Thread_Worker(storage.Task_Chan)
	fmt.Println("Started Worker")

	go storage.Cleaner()
	fmt.Println("Started Cleaner")

	monitor_conn, err := connection.Connect_Monitor()
	if err != nil {
		fmt.Println("Monitor Connection Error")
		os.Exit(1)
	}
	go connection.HandleMonitorConnection(monitor_conn)

	if replication.Metadata.Role == "slave" {
		master_conn, err := replication.MakeHandShake()
		if err != nil {
			os.Exit(1)
		}
		go connection.HandleMasterConnection(master_conn)
	}

	if replication.Metadata.Role == "master" {
		go storage.Dumper()
	}

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go connection.HandleConnection(con)
	}
}
