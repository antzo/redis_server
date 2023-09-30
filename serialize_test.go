package redis_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialize(t *testing.T) {
	testCases := []struct {
		desc  string
		input Command
		want  struct {
			c   string
			err error
		}
	}{
		{
			desc: "hello world SimpleString",
			input: Command{
				Type: "SimpleString",
				Data: "hello world",
			},
			want: struct {
				c   string
				err error
			}{
				c:   "+hello world\r\n",
				err: nil,
			},
		},
		{
			desc:  "null value BulkString",
			input: Command{Type: BulkStringType},
			want: struct {
				c   string
				err error
			}{
				c:   "$-1\r\n",
				err: nil,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := Serialize(tC.input)

			assert.Equal(t, tC.want.c, got)
			assert.Equal(t, tC.want.err, err)
		})
	}
}
