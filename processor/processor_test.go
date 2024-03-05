package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeInput(t *testing.T) {
	tests := []struct {
		name              string
		s                 string
		lenIncompleteLine int
		line              string
		exitLoop          bool
		emptyString       bool
	}{
		{
			name:              "successfully sanitizes input when line is complete",
			s:                 "> 1544206563 Delta +0.5",
			lenIncompleteLine: 0,
			line:              "1544206563 Delta +0.5",
			exitLoop:          false,
			emptyString:       false,
		},
		{
			name:              "successfully sanitizes input when line is incomplete",
			s:                 "> 1544206563 Delta",
			lenIncompleteLine: 0,
			line:              "1544206563 Delta",
			exitLoop:          false,
			emptyString:       false,
		},
		{
			name:              "successfully sanitizes input when EOF included ",
			s:                 "> 1544206563 Delta EOF",
			lenIncompleteLine: 0,
			line:              "1544206563 Delta ",
			exitLoop:          true,
			emptyString:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, el, es := sanitizeInput(tt.s, tt.lenIncompleteLine)
			assert.Equal(t, l, tt.line)
			assert.Equal(t, el, tt.exitLoop)
			assert.Equal(t, es, tt.emptyString)
		})
	}
}

func Test_isEmpty(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: true,
		},
		{
			name: "empty string",
			s:    "not empty",
			want: false,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, isEmpty(tt.s), tt.want)
	}
}
