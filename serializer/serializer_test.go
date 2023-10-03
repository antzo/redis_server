package serializer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Deserialize(t *testing.T) {
	testCases := []struct {
		desc string
		want struct {
			message Message
			err     error
		}
		binaryData []byte
	}{
		{
			desc:       "deserialize SimpleString",
			binaryData: []byte("+hello world\r\n"),
			want: struct {
				message Message
				err     error
			}{
				message: StandardMessage{
					typeName: simpleStringType,
					data:     []byte("+hello world\r\n"),
				},
				err: nil,
			},
		},
		{
			desc:       "deserialize empty SimpleString",
			binaryData: []byte("+\r\n"),
			want: struct {
				message Message
				err     error
			}{
				message: StandardMessage{
					typeName: simpleStringType,
					data:     []byte("+\r\n"),
				},
				err: nil,
			},
		},
		{
			desc:       "simple Integer",
			binaryData: []byte(":123"),
			want: struct {
				message Message
				err     error
			}{
				message: StandardMessage{
					typeName: integerType,
					data:     []byte(":123"),
				},
				err: nil,
			},
		},
		{
			desc:       "bulkString",
			binaryData: []byte("$5\r\nhello\r\n"),
			want: struct {
				message Message
				err     error
			}{
				message: StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$5\r\nhello\r\n"),
				},
				err: nil,
			},
		},
		{
			desc: "deserialize array of 2 elements with bulk strings",
			want: struct {
				message Message
				err     error
			}{
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
				err: nil,
			},
			binaryData: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
		},
		{
			desc: "deserialize array of arrays with a bulk string",
			want: struct {
				message Message
				err     error
			}{
				message: ArrayMessage{Messages: []Message{
					ArrayMessage{Messages: []Message{
						StandardMessage{
							typeName: bulkStringType,
							data:     []byte("$5\r\nhello\r\n"),
						},
					}},
				}},
				err: nil,
			},
			binaryData: []byte("*1\r\n*1\r\n$5\r\nhello\r\n"),
		},
		{
			desc: "initial redis-cli message",
			want: struct {
				message Message
				err     error
			}{
				message: ArrayMessage{Messages: []Message{
					StandardMessage{typeName: bulkStringType, data: []byte("$7\r\nCOMMAND\r\n")},
					StandardMessage{typeName: bulkStringType, data: []byte("$4\r\nDOCS\r\n")},
				}},
				err: nil,
			},
			binaryData: []byte("*2\r\n$7\r\nCOMMAND\r\n$4\r\nDOCS\r\n"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := Deserialize(tC.binaryData)

			assert.Equal(t, tC.want.message, got)
			assert.Equal(t, tC.want.err, err)
		})
	}
}
