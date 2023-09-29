package redis_server

import (
	"errors"
	"fmt"
	"strings"
)

const separator = "\r\n"

type Command struct {
	Type string
	Data string
}

func Deserialize(commandString string) ([]Command, error) {
	cmd := make([]Command, 0)

	if len(commandString) <= 0 {
		return cmd, errors.New("empty string")
	}

	commandType, err := parseCommandType(commandString[0])
	if err != nil {
		return cmd, err
	}

	return parseCommands(commandType, commandString)
}

func parseCommands(commandType, commandString string) ([]Command, error) {
	cmd := make([]Command, 0)

	switch commandType {
	case "SimpleString":
		crlf := strings.Index(commandString, separator)
		if crlf == -1 {
			return cmd, errors.New("command not terminated with CRLF")
		}
		c := Command{}
		c.Type = commandType
		c.Data = commandString[1:crlf]
		cmd = append(cmd, c)
	case "BulkStrings":
		parsedCmd, e := parseBulkString(commandString)
		if e != nil {
			return cmd, e
		}

		for _, v := range parsedCmd {
			cmd = append(cmd, v)
		}
	case "Arrays":
		first := strings.Index(commandString, separator)
		startNextCommands := commandString[first+len(separator):]
		nextType, err := parseCommandType(startNextCommands[0])
		if err != nil {
			return cmd, err
		}
		return parseCommands(nextType, startNextCommands)

	case "SimpleErrors":
	case "Integers":
	default:
		return cmd, errors.New("not implemented")
	}

	return cmd, nil
}

func parseCommandType(firstByte uint8) (string, error) {
	switch firstByte {
	case '+':
		return "SimpleString", nil
	case '-':
		return "SimpleErrors", nil
	case ':':
		return "Integers", nil
	case '$':
		return "BulkStrings", nil
	case '*':
		return "Arrays", nil
	default:
		return "", fmt.Errorf("%c is unsupported command data type", firstByte)
	}
}

func parseBulkString(c string) ([]Command, error) {
	res := make([]Command, 0)
	cmds := strings.Split(c, separator)

	for i := 1; i < len(cmds); i = i + 2 {
		if cmds[i] == "" {
			continue
		}
		cmd := Command{
			Type: "BulkStrings",
			Data: cmds[i],
		}
		res = append(res, cmd)
	}

	return res, nil
}
