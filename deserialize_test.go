package redis_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeserialize(t *testing.T) {
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
					{Type: SimpleStringType, Data: "hello world"},
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
				c: []Command{
					{Type: BulkStringType},
				},
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
					{Type: BulkStringType, Data: "echo"},
					{Type: BulkStringType, Data: "hello world"},
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
					{Type: BulkStringType, Data: "ping"},
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
					{Type: BulkStringType, Data: "get"},
					{Type: BulkStringType, Data: "key"},
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
					{Type: SimpleStringType, Data: "OK"},
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
				c: []Command{
					{Type: BulkStringType},
				},
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
