package storage

import (
	"strings"
	"sync"
	"time"
)

var Local_Storage_Val = make(map[string]string)
var Local_Storage_Time = make(map[string]int64)
var Write_Mutex sync.RWMutex

func SetValue(key string, val string, time int64) {
	Write_Mutex.Lock()
	Local_Storage_Val[key] = val
	Local_Storage_Time[key] = time
	Write_Mutex.Unlock()
}

func GetValue(key string) string {
	Write_Mutex.RLock()
	val, err1 := Local_Storage_Val[key]
	ti, err := Local_Storage_Time[key]
	Write_Mutex.RUnlock()
	if !err || !err1 {
		val = "-1"
	} else if ti != -1 && time.Now().UnixMilli() > ti {
		val = "-1"
	}
	return val
}

func HasPrefix(key string) string {
	k := ""
	for _, c := range key {
		if c == '*' {
			break
		}
		k = k + string(c)
	}
	var v []string = make([]string, 0)
	for key, val := range Local_Storage_Val {
		if strings.HasPrefix(key, k) {
			v = append(v, key+"\t"+val)
		}
	}
	return strings.Join(v, "\n")
}
