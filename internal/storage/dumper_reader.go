package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/config"
)

const DumperTime = 10

func Dumper() {
	for {
		time.Sleep(DumperTime * time.Second)
		file, err := os.OpenFile("./"+config.RedisConfig.Dir+"/"+config.RedisConfig.DbFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(file)
		for key, val := range Local_Storage_Val {
			writer.WriteString(key + "$" + val + "$" + strconv.FormatInt(Local_Storage_Time[key], 10) + "\n")
		}
		SlaveOffsetVal = MasterReplOffset
		writer.Flush()
		file.Close()
	}
}

func Reader() {
	file, err := os.Open("./" + config.RedisConfig.Dir + "/" + config.RedisConfig.DbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	Local_Storage_Val = make(map[string]string)
	Local_Storage_Time = make(map[string]int64)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "$")
		fmt.Println(parts)
		if len(parts) == 3 {
			Local_Storage_Val[parts[0]] = parts[1]
			ti, _ := strconv.ParseInt(parts[2], 10, 64)
			Local_Storage_Time[parts[0]] = ti
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
