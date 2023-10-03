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
}

var availableCommands = []string{"ping", "echo"}

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
		cmdName := strings.ToLower(string(d[nextLF : nextLF+msgLength]))

		for _, v := range availableCommands {
			if v == cmdName {
				cmd.name = cmdName
				return cmd, nil
			}
		}
	}

	return cmd, errors.New("no command found")
}
