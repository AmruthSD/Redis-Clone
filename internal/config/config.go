package config

import "flag"

type Config struct {
	Dir        string
	DbFileName string
}

var RedisConfig Config

func LoadConfig() {
	flag.StringVar(&RedisConfig.Dir, "dir", "", "Directory to store RDB file")
	flag.StringVar(&RedisConfig.DbFileName, "dbfilename", "dump.rdb", "RDB file name")

	flag.Parse()
}
