package replication

import (
	"math/rand"
	"time"
)

const alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func new_replication_id() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 40; i++ {
		s = s + string(alphanumeric[rand.Int()%len(alphanumeric)])
	}
	return s
}
