package redis_server

import (
	"errors"
	"fmt"
)

func Serialize(cmd Command) (string, error) {
	switch cmd.Type {
	case SimpleStringType:
		return "+" + cmd.Data + separator, nil
	case BulkStringType:
		l := len(cmd.Data)
		if l == 0 {
			l = -1
		}
		return "$" + fmt.Sprintf("%d", l) + cmd.Data + separator, nil
	default:
		return "", errors.New("not implemented")
	}
}
