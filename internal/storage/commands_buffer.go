package storage

type Command struct {
	Parts     []string
	Num_Bytes int
}

const LIMIT int = 10

func InsertCommand(parts []string, num_bytes int) {
	node := Command{Parts: parts, Num_Bytes: num_bytes}
	for {
		if Prev_Commands.Len() >= LIMIT {
			n := Prev_Commands.Front()
			Prev_Commands.Remove(n)
		} else {
			break
		}
	}
	Prev_Commands.PushBack(node)
}
