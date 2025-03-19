package replication

import "strings"

func SendMessageToSlaves(parts []string) {
	msg := strings.Join(parts, " ")
	if Metadata.Role == "master" {
		for key, val := range SlavesConnections {
			if val {
				ch, ex := ConnectionChannels[key]
				if ex {
					ch <- msg
				}
			}
		}
	}
}
