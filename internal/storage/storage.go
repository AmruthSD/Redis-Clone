package storage

import (
	"sync"
	"time"
)

var local_storage_val = make(map[string]string)
var local_storage_time = make(map[string]int64)
var write_mutex sync.RWMutex

func SetValue(key string, val string, time int64) {
	write_mutex.Lock()
	local_storage_val[key] = val
	local_storage_time[key] = time
	write_mutex.Unlock()
}

func GetValue(key string) string {
	write_mutex.RLock()
	val, err1 := local_storage_val[key]
	ti, err := local_storage_time[key]
	write_mutex.RUnlock()
	if !err || !err1 {
		val = "-1"
	} else if ti != -1 && time.Now().UnixMilli() > ti {
		val = "-1"
	}
	return val
}
