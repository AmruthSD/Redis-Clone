package storage

import "errors"

type Command struct {
	Parts     []string
	Num_Bytes int
}

const LIMIT int = 10

var TotalBytes = 0

func InsertCommand(parts []string, num_bytes int) {
	node := Command{Parts: parts, Num_Bytes: num_bytes}
	for {
		if Prev_Commands.Len() >= LIMIT {
			n := Prev_Commands.Front()
			val := n.Value.(Command)
			TotalBytes -= val.Num_Bytes
			Prev_Commands.Remove(n)
		} else {
			break
		}
	}
	Prev_Commands.PushBack(node)
	TotalBytes += num_bytes
}

func LastFewCommands(num_bytes int) ([][]string, error) {
	if num_bytes > TotalBytes {
		return nil, errors.New("not enough logs")
	}
	var buf [][]string
	num := 0
	node := Prev_Commands.Back()
	for num < num_bytes {
		if num+node.Value.(Command).Num_Bytes > num_bytes {
			return nil, errors.New("out of sync")
		}
		num += node.Value.(Command).Num_Bytes
		buf = append(buf, node.Value.(Command).Parts)
		node = node.Prev()
	}
	return buf, nil
}
