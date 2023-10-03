package redis_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	expected := StandardMessage{
		typeName: simpleStringType,
		data:     []byte("+PONG\r\n"),
	}

	assert.Equal(t, expected, Ping(StandardMessage{
		typeName: simpleStringType,
		data:     []byte("+PING\r\n"),
	}))
}

func TestGetCommand(t *testing.T) {
	testCases := []struct {
		desc  string
		input Message
		want  struct {
			cmd Command
			err error
		}
	}{
		{
			desc: "ping command",
			input: ArrayMessage{Messages: []Message{
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$4\r\nping\r\n"),
				},
			}},
			want: struct {
				cmd Command
				err error
			}{
				cmd: Command{name: "ping"},
				err: nil,
			},
		},
		{
			desc: "ping command in uppercase",
			input: ArrayMessage{Messages: []Message{
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$4\r\nPING\r\n"),
				},
			}},
			want: struct {
				cmd Command
				err error
			}{
				cmd: Command{name: "ping"},
				err: nil,
			},
		},
		{
			desc: "ECHO hello world",
			input: ArrayMessage{Messages: []Message{
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$4\r\nECHO\r\n"),
				},
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$5\r\nhello\r\n"),
				},
				StandardMessage{
					typeName: bulkStringType,
					data:     []byte("$5\r\nworld\r\n"),
				},
			}},
			want: struct {
				cmd Command
				err error
			}{
				cmd: Command{name: "echo", args: [][]byte{
					[]byte("hello\r\n"),
					[]byte("world\r\n"),
				}},
				err: nil,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := GetCommand(tC.input)

			assert.Equal(t, tC.want.cmd, got)
			assert.Equal(t, tC.want.err, err)
		})
	}
}
