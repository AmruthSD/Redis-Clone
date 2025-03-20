package storage

import (
	"container/list"
	"strings"
	"time"
)

var Local_Storage_Val = make(map[string]string)
var Local_Storage_Time = make(map[string]int64)

var Prev_Commands = list.New()

type Task struct {
	Fn        func() any
	Result_ch chan any
}

var Task_Chan = make(chan Task)

func Single_Thread_Worker(tasks <-chan Task) {
	for task := range tasks {
		res := task.Fn()
		if res != nil {
			task.Result_ch <- res
		}
	}
}

func SetValue(key string, val string, time int64) {

	Local_Storage_Val[key] = val
	Local_Storage_Time[key] = time

}

func GetValue(key string) string {

	val, err1 := Local_Storage_Val[key]
	ti, err := Local_Storage_Time[key]

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
