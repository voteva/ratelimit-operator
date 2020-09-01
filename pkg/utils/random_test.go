package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_BuildRandomString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	str := BuildRandomString(5)
	a.Equal(len(str), 5)
	for i := 0; i < len(str); i++ {
		a.GreaterOrEqual(strings.IndexByte(charset, str[i]), 0)
	}
}

func Test_BuildRandomInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	n := BuildRandomInt(10)
	a.GreaterOrEqual(n, 0)
	a.Less(n, 10)
}
