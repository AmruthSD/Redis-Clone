package storage

var local_storage = make(map[string]string)

func SetValue(key string, val string) {
	local_storage[key] = val
}

func GetValue(key string) string {
	val, err := local_storage[key]
	if !err {
		val = "NULL"
	}
	return val
}
