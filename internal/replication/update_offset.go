package replication

func UpdateOffset(parts []string) {
	tt := 0
	for _, v := range parts {
		tt += len(v)
	}
	Metadata.MasterReplOffset += tt
}
