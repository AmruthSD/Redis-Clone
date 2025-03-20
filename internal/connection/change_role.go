package connection

import (
	"os"

	"github.com/AmruthSD/Redis-Clone/internal/replication"
	"github.com/AmruthSD/Redis-Clone/internal/storage"
)

func Slave_Init() {
	master_conn, err := replication.MakeHandShake()
	if err != nil {
		os.Exit(1)
	}
	go HandleMasterConnection(master_conn)
}

func Master_Init() {
	go storage.Dumper()
}
