package connection

import (
	"errors"
	"net"
)

var parse_func = map[string]func([]string, net.Conn){
	"PING":     handle_ping,
	"ECHO":     handle_echo,
	"SET":      handle_set,
	"GET":      handle_get,
	"CONFIG":   handle_config,
	"KEYS":     handle_keys,
	"INFO":     handle_info,
	"REPLCONF": handle_replconf,
	"PSYNC":    handle_psync,
	"OK":       handle_ok,
	"DEL":      handle_del,
	"unknown":  handle_unknown,
}

func Parse(parts []string) (func([]string, net.Conn), error) {
	if len(parts) == 0 {
		return nil, errors.New("empty command")
	}
	if handler, exists := parse_func[parts[0]]; exists {
		return handler, nil
	} else {
		return nil, errors.New("unknown command")
	}
}
