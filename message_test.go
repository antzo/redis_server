package redis_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayMessage_Data(t *testing.T) {
	testCases := []struct {
		desc    string
		message ArrayMessage
		want    []byte
	}{
		{
			desc: "array with two bulkString elements of 5 length",
			want: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			message: ArrayMessage{Messages: []Message{
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$5\r\nhello\r\n"),
				},
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$5\r\nworld\r\n"),
				},
			}},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, string(tC.want), string(tC.message.Data()))
		})
	}
}
