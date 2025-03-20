package replication

type MasterSlaveData struct {
	Role             string
	NumberOfSlaves   int
	MasterReplid     string
	MasterReplOffset int
	MasterAddress    string
}

func NewMasterSlaveData() MasterSlaveData {
	return MasterSlaveData{
		Role:             "",
		NumberOfSlaves:   0,
		MasterReplid:     new_replication_id(),
		MasterAddress:    "",
		MasterReplOffset: 0,
	}
}

var Metadata MasterSlaveData = NewMasterSlaveData()

var SlavesConnections map[string]bool = make(map[string]bool)

var ConnectionChannels map[string]chan string = make(map[string]chan string)
