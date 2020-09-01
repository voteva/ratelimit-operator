package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Merge(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	m1 := map[string]string{"a": "1", "b": "2"}
	m2 := map[string]string{"c": "3"}

	m3 := Merge(m1, m2)

	a.Equal(m1, map[string]string{"a": "1", "b": "2"})
	a.Equal(m2, map[string]string{"c": "3"})
	a.Equal(m3, map[string]string{"a": "1", "b": "2", "c": "3"})

}
