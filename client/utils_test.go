package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatName(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{name: "name", title: "abcd", want: "abcd"},
		{name: "name - 2", title: "abcd|bcded", want: "abcd"},
		{"name - 3", "abcd&bcd|ed", "abcdbcd"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, formatName(tt.title), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
