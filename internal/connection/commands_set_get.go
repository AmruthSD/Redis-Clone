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
		storage.SetValue(parts[1], parts[2], -1)
	} else if len(parts) == 5 && parts[3] == "PX" {
		ext, _ := strconv.ParseInt(parts[4], 10, 64)
		ti := time.Now().UnixMilli() + ext
		storage.SetValue(parts[1], parts[2], ti)
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
	val := storage.GetValue(parts[1])
	conn.Write([]byte(val))
}
