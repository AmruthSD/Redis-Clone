package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Dir        string
	DbFileName string
	Port       int
	ReplicaOf  string
}

var RedisConfig Config

func LoadConfig() {
	flag.StringVar(&RedisConfig.Dir, "dir", "persist", "Directory to store RDB file")
	flag.StringVar(&RedisConfig.DbFileName, "dbfilename", "dump.rdb", "RDB file name")
	flag.IntVar(&RedisConfig.Port, "port", 6379, "Port to which you want the application to listen to")
	flag.StringVar(&RedisConfig.ReplicaOf, "replicaof", "", "Port of the master")

	fmt.Println("Flags read")

	flag.Parse()
}
