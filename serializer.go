package redis_server

import (
	"bytes"
	"errors"
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

	return nil, nil
}

func deserializeBulkString(data []byte) (Message, error) {
	message := StandardMessage{typeName: bulkStringType}

	length, err := strconv.Atoi(string(data[1]))
	if err != nil {
		return message, err
	}

	if len(data) < length+6 {
		return message, errors.New("parser: deserialize: bulkString: invalid message")
	}

	message.data = data[0 : length+6]

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
		command.Messages = append(command.Messages, cmd)
		next = next + len(cmd.Data())
	}

	return command, nil
}
