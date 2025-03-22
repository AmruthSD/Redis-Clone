package main

import (
	"fmt"
	"net"
)

func SendSlave(conn net.Conn) {
	fmt.Println(ConnctionMap)
	fmt.Println(ConnectionCount)
	min_adress := ""
	min_cnt := 100000
	for key, val := range ConnectionCount {
		if key != MasterKey && val < int(min_cnt) && ConnctionMap[key] != "-1" {
			min_adress = key
		}
	}
	if min_adress == "" {
		min_adress = MasterKey
	}
	ConnectionCount[min_adress]++
	conn.Write([]byte(ConnctionMap[min_adress]))
	fmt.Println(ConnctionMap[min_adress])

}
