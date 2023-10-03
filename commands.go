package redis_server

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MessageHandler func(message Message) Message

func Ping(_ Message) Message {
	return StandardMessage{
		typeName: simpleStringType,
		data:     []byte(fmt.Sprintf("+PONG%c%c", CR, LF)),
	}
}

type Command struct {
	name string
	args [][]byte
}

var availableCommands = []string{"ping", "echo"}

var errCommandNotFound = errors.New("no command found")

func GetCommand(message Message) (Command, error) {
	cmd := Command{}

	if cmdArr, ok := message.(ArrayMessage); ok && len(cmdArr.Messages) >= 0 {
		d := cmdArr.Messages[0].Data()
		nextLF := bytes.IndexByte(d, LF)

		msgLength, err := strconv.Atoi(string(d[1 : nextLF-1]))
		if err != nil {
			return cmd, err
		}
		nextLF++
		cmd.name = strings.ToLower(string(d[nextLF : nextLF+msgLength]))

		if !commandExist(cmd.name) {
			return cmd, errCommandNotFound
		}

		if len(cmdArr.Messages) == 1 {
			return cmd, nil
		}

		// Arguments
		args := make([][]byte, 0)
		for i := 1; i < len(cmdArr.Messages); i++ {
			d = cmdArr.Messages[i].Data()
			nextLF = bytes.IndexByte(d, LF)
			msgLength, err = strconv.Atoi(string(d[1 : nextLF-1]))
			if err != nil {
				return cmd, err
			}
			nextLF++
			args = append(args, d[nextLF:])
		}
		cmd.args = args

		return cmd, nil
	}

	return cmd, errCommandNotFound
}

func commandExist(cmdName string) bool {
	for _, v := range availableCommands {
		if v == cmdName {
			return true
		}
	}
	return false
}
