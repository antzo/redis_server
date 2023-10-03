package serializer

import (
	"bytes"
	"fmt"
)

const (
	simpleStringType = '+'
	integerType      = ':'
	simpleErrorType  = '-'
	bulkStringType   = '$'
	arrayType        = '*'
)

// Message defines a RESP2 protocol message
type Message interface {
	Type() rune
	Data() []byte
}

// StandardMessage is used to represent BulkStrings, Integers, SimpleStrings and SimpleErrors
type StandardMessage struct {
	typeName rune
	data     []byte
}

func (c StandardMessage) Type() rune   { return c.typeName }
func (c StandardMessage) Data() []byte { return c.data }

// ArrayMessage represents an Array RESP2 protocol request/response
type ArrayMessage struct {
	Messages []Message
}

func (c ArrayMessage) Type() rune { return arrayType }

// Data writes RESP2 array binary data and all the commands inside the array
func (c ArrayMessage) Data() []byte {
	var buf bytes.Buffer

	buf.WriteByte(arrayType)
	buf.WriteString(fmt.Sprintf("%d", len(c.Messages)))
	buf.WriteByte(CR)
	buf.WriteByte(LF)

	for _, message := range c.Messages {
		buf.Write(message.Data())
	}

	return buf.Bytes()
}
