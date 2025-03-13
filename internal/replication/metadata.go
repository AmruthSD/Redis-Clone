package replication

type MasterSlaveData struct {
	Role             string
	NumberOfSlaves   int
	MasterReplid     string
	MasterReplOffset int
}

func NewMasterSlaveData() MasterSlaveData {
	return MasterSlaveData{
		Role:             "master",
		NumberOfSlaves:   0,
		MasterReplid:     new_replication_id(),
		MasterReplOffset: 0,
	}
}

var Metadata MasterSlaveData = NewMasterSlaveData()
