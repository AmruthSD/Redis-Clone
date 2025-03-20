package replication

import "github.com/AmruthSD/Redis-Clone/internal/storage"

func UpdateOffset(parts []string) {
	tt := 0
	for _, v := range parts {
		tt += len(v)
	}
	Metadata.MasterReplOffset += tt
	if Metadata.Role == "master" {
		storage.InsertCommand(parts, tt)
	}
}
