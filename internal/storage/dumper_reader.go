package storage

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Dumper() {
	for {
		time.Sleep(900 * time.Second)
		file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(file)
		for key, val := range Local_Storage_Val {
			writer.WriteString(key + "$" + val + "$" + strconv.FormatInt(Local_Storage_Time[key], 10) + "\n")
		}
		writer.Flush()
		file.Close()
	}
}

func Reader() {
	file, err := os.Open("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	Local_Storage_Val = make(map[string]string)
	Local_Storage_Time = make(map[string]int64)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "&")
		Local_Storage_Val[parts[0]] = parts[1]
		ti, _ := strconv.ParseInt(parts[2], 10, 64)
		Local_Storage_Time[parts[0]] = ti
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
