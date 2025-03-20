package main

import (
	"fmt"
	"os"
	"time"
)

const Timeout = 10

func CheckMasterAlive() {
	time.Sleep(Timeout * time.Second)
	ex := ConnctionMap[MasterKey]
	if ex {
		return
	}

	MasterKey = ""
	for key := range ConnctionMap {
		MasterKey = key
		break
	}
	if MasterKey == "" {
		fmt.Println("All Master Candidates Closed")
		os.Exit(2)
	}
	for key := range ConnctionMap {
		go func() {
			if key == MasterKey {
				ConnectionChan[key] <- "YOU ARE THE MASTER"
			} else {
				ConnectionChan[key] <- MasterKey
			}
		}()
	}
}
