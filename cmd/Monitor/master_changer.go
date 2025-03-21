package main

import (
	"fmt"
	"os"
	"time"
)

const Timeout = 1

func CheckMasterAlive() {
	for {
		time.Sleep(Timeout * time.Second)
		fmt.Println("Checking Master Status...")

		_, ex := ConnctionMap[MasterKey]
		if ex {
			continue
		}

		fmt.Println("Master Closed")

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
}
