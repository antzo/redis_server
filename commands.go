package redis_server

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MessageHandler func(command Command) Message

// Ping returns PONG if no argument is provided,
// otherwise return a copy of the argument as a bulk. This command is useful for:
// Testing whether a connection is still alive.
// Verifying the server's ability to serve data - an error is returned when this isn't
// the case (e.g., during load from persistence or accessing a stale replica).
// Measuring latency.
// If the client is subscribed to a channel or a pattern, it will instead return a multi-bulk
// with a "pong" in the first position and an empty bulk in the second position,
// unless an argument is provided in which case it returns a copy of the argument.
//
// Return
// - Simple string reply, and specifically PONG, when no argument is provided.
// - Bulk string reply the argument provided, when applicable.
func Ping(cmd Command) Message {
	if len(cmd.args) == 0 {
		return StandardMessage{
			typeName: simpleStringType,
			data:     []byte(fmt.Sprintf("%cPONG%c%c", simpleStringType, CR, LF)),
		}
	}

	args := bytes.Join(cmd.args, []byte(" "))
	bulk := []byte(
		fmt.Sprintf("%c%d%c%c%s%c%c", bulkStringType, len(args), CR, LF, args, CR, LF),
	)

	return StandardMessage{
		typeName: bulkStringType,
		data:     bulk,
	}
}

// Echo returns message of bulkStringType
func Echo(command Command) Message {
	var buf bytes.Buffer

	d := bytes.Join(command.args, []byte(" "))
	buf.WriteByte(bulkStringType)
	buf.Write([]byte(fmt.Sprintf("%d%c%c", len(d), CR, LF)))
	buf.Write([]byte(fmt.Sprintf("%s%c%c", d, CR, LF)))

	return StandardMessage{
		typeName: bulkStringType,
		data:     buf.Bytes(),
	}
}

type Command struct {
	name string
	args [][]byte
}

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

		if len(cmdArr.Messages) == 1 {
			return cmd, nil
		}

		// Arguments
		args := make([][]byte, 0)
		for i := 1; i < len(cmdArr.Messages); i++ {
			d = cmdArr.Messages[i].Data()
			nextLF = bytes.IndexByte(d, LF)

			// Get message length
			msgLength, err = strconv.Atoi(string(d[1 : nextLF-1]))
			if err != nil {
				return cmd, err
			}

			// Add bulkString without CRLF
			nextLF++
			d = d[nextLF:]
			args = append(args, d[:msgLength])
		}
		cmd.args = args

		return cmd, nil
	}

	return cmd, errCommandNotFound
}
