package connection

import (
	"net"
	"strconv"
	"time"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func handle_set(parts []string, conn net.Conn) {
	go replication.SendMessageToSlaves(parts)
	if len(parts) == 3 {
		replication.UpdateOffset(parts)

		storage.Task_Chan <- storage.Task{Fn: func() any { storage.SetValue(parts[1], parts[2], -1); return nil }, Result_ch: nil}
	} else if len(parts) == 5 && parts[3] == "PX" {
		ext, _ := strconv.ParseInt(parts[4], 10, 64)
		ti := time.Now().UnixMilli() + ext
		replication.UpdateOffset(parts)
		storage.Task_Chan <- storage.Task{Fn: func() any { storage.SetValue(parts[1], parts[2], ti); return nil }, Result_ch: nil}
	} else if len(parts) != 3 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	}
	conn.Write([]byte("OK"))
}

func handle_get(parts []string, conn net.Conn) {
	if len(parts) != 2 {
		conn.Write([]byte("Argument Count Not Right"))
		return
	}
	result_ch := make(chan any)
	storage.Task_Chan <- storage.Task{Fn: func() any { val := storage.GetValue(parts[1]); return val }, Result_ch: result_ch}
	conn.Write(to_bytes(<-result_ch))
}

func to_bytes(value any) []byte {
	switch v := value.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}
