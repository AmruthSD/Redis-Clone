package storage

import (
	"encoding/binary"
	"os"
)

const (
	opCodeModuleAux    byte = 247 /* Module auxiliary data. */
	opCodeIdle         byte = 248 /* LRU idle time. */
	opCodeFreq         byte = 249 /* LFU frequency. */
	opCodeAux          byte = 250 /* RDB aux field. */
	opCodeResizeDB     byte = 251 /* Hash table resize hint. */
	opCodeExpireTimeMs byte = 252 /* Expire time in milliseconds. */
	opCodeExpireTime   byte = 253 /* Old expire time in seconds. */
	opCodeSelectDB     byte = 254 /* DB number of the following keys. */
	opCodeEOF          byte = 255
)

func sliceIndex(data []byte, sep byte) int {
	for i, b := range data {
		if b == sep {
			return i
		}
	}
	return -1
}

func parseTable(bytes []byte) []byte {
	start := sliceIndex(bytes, opCodeResizeDB)
	end := sliceIndex(bytes, opCodeEOF)
	return bytes[start+1 : end]
}

func ReadFile(path string) {
	c, _ := os.ReadFile(path)
	key := parseTable(c)
	records, _ := parseRecords(key)
	for _, rec := range records {
		SetValue(rec.Key, rec.Value, rec.Expiry)
	}
}

type Record struct {
	Key    string
	Value  string
	Expiry int64
}

func parseRecords(data []byte) ([]Record, error) {
	var records []Record
	i := 0
	for i < len(data) {
		kLen := int(data[i])
		i++
		key := string(data[i : i+kLen])
		i += kLen
		vLen := int(data[i])
		i++
		value := string(data[i : i+vLen])
		i += vLen
		expiryFlag := data[i]
		i++

		expiry := int64(-1)
		if expiryFlag == 1 {
			expiry = int64(binary.BigEndian.Uint64(data[i : i+8]))
			i += 8
		}
		records = append(records, Record{
			Key:    key,
			Value:  value,
			Expiry: expiry,
		})
	}
	return records, nil
}
