package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/westwin/gocommon/util"
)

func TestRandString(t *testing.T) {
	len := 10
	choices := `abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!"#$%&()*+:;=<>=?@][~]{}}`

	s1 := util.RandString(len, choices)
	assert.Len(t, s1, len)
	t.Logf("rand string 1: %s", s1)

	s2 := util.RandString(len, choices)
	assert.Len(t, s2, len)
	t.Logf("rand string 1: %s", s2)
	assert.NotEqual(t, s1, s2) // the chance is very small
}
