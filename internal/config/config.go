package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Dir        string
	DbFileName string
	Port       int
}

var RedisConfig Config

func LoadConfig() {
	flag.StringVar(&RedisConfig.Dir, "dir", "db", "Directory to store RDB file")
	flag.StringVar(&RedisConfig.DbFileName, "dbfile", "dump.rdb", "RDB file name")
	flag.IntVar(&RedisConfig.Port, "port", 6379, "Port to which you want the application to listen to")

	fmt.Println("Flags read")

	flag.Parse()
}
