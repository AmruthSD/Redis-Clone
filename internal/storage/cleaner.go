package storage

import "time"

func Cleaner() {
	for {
		time.Sleep(10 * time.Second)
		go func() {
			for key := range Local_Storage_Val {
				if Local_Storage_Time[key] != -1 && Local_Storage_Time[key] > time.Now().UnixMilli() {
					delete(Local_Storage_Time, key)
					delete(Local_Storage_Val, key)
				}
			}
		}()
	}
}
