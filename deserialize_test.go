package redis_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialize(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  struct {
			c   []Command
			err error
		}
	}{
		{
			desc:  "hello world SimpleString",
			input: "+hello world\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c: []Command{
					{Type: "SimpleString", Data: "hello world"},
				},
				err: nil,
			},
		},
		{
			desc:  "null value BulkString",
			input: "$-1\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c:   []Command{},
				err: nil,
			},
		},
		{
			desc:  "multiple BulkStrings",
			input: "*2\r\n$4\r\necho\r\n$5\r\nhello world\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c: []Command{
					{Type: "BulkStrings", Data: "echo"},
					{Type: "BulkStrings", Data: "hello world"},
				},
				err: nil,
			},
		},
		{
			desc:  "ping Arrays with BulkString",
			input: "*1\r\n$4\r\nping\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c: []Command{
					{Type: "BulkStrings", Data: "ping"},
				},
				err: nil,
			},
		},
		{
			desc:  "get/key Array",
			input: "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c: []Command{
					{Type: "BulkStrings", Data: "get"},
					{Type: "BulkStrings", Data: "key"},
				},
				err: nil,
			},
		},
		{
			desc:  "OK SimpleString",
			input: "+OK\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c: []Command{
					{Type: "SimpleString", Data: "OK"},
				},
				err: nil,
			},
		},
		{
			desc:  "empty array",
			input: "$0\r\n\r\n",
			want: struct {
				c   []Command
				err error
			}{
				c:   []Command{},
				err: nil,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, err := Deserialize(tC.input)

			assert.Equal(t, tC.want.c, res)
			assert.Equal(t, tC.want.err, err)
		})
	}
}
