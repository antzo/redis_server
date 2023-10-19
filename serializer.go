package redis_server

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const (
	CR = '\r'
	LF = '\n'
)

func Deserialize(binaryData []byte) (Message, error) {
	if len(binaryData) <= 0 {
		return nil, nil
	}

	messageType := binaryData[0]

	switch messageType {
	case arrayType:
		return deserializeArray(binaryData)
	case simpleStringType:
		return StandardMessage{typeName: simpleStringType, data: binaryData}, nil
	case integerType:
		return StandardMessage{typeName: integerType, data: binaryData}, nil
	case bulkStringType:
		return deserializeBulkString(binaryData)
	}

	return nil, fmt.Errorf("messageType: %c not supported on data: %s", messageType, binaryData)
}

func deserializeBulkString(data []byte) (Message, error) {
	message := StandardMessage{typeName: bulkStringType}
	next := bytes.IndexByte(data, CR)

	length, err := strconv.Atoi(string(data[1:next]))
	if err != nil {
		return message, err
	}

	if len(data) < length+6 {
		return message, errors.New("parser: deserialize: bulkString: invalid message")
	}

	// This will add how many bytes we need to read.
	// For doing that we need to add 6+one byte per digit on the length.
	// - one byte for $ char that indicates it's a bulkstring (1)
	// - two bytes on the first \r\n (3)
	// - two more for the final \r\n (5)
	// - and finally N bytes for each digit that indicates the length (ex: 3length -> 1 byte, 300length -> 3bytes)
	digits := length / 10
	message.data = data[0 : length+6+digits]

	return message, nil
}

func deserializeArray(data []byte) (Message, error) {
	if len(data) <= 1 {
		return nil, errors.New("parser: deserializer: array message without length received")
	}

	numCommands, err := strconv.Atoi(string(data[1]))
	if err != nil {
		return nil, err
	}

	command := ArrayMessage{}
	next := bytes.IndexByte(data, LF)

	for i := 0; i < numCommands; i++ {
		cmd, err := Deserialize(data[next+1:])
		if err != nil {
			return nil, err
		}
		if cmd == nil {
			return nil, nil
		}
		command.Messages = append(command.Messages, cmd)
		next = next + len(cmd.Data())
	}

	return command, nil
}
