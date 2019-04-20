package util_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/westwin/gocommon/util"
)

func TestSliceContains(t *testing.T) {
	slice := []string{"first", "second", "third"}

	for _, ele := range slice {
		assert.True(t, util.Contains(slice, ele))
	}

	assert.False(t, util.Contains(slice, "not-existed"))

	assert.False(t, util.Contains([]string{}, "empty slice"))
}

func TestSliceContainsIgnoreCase(t *testing.T) {
	slice := []string{"first", "second", "third"}

	for _, ele := range slice {
		assert.True(t, util.ContainsIgnoreCase(slice, ele))
		assert.True(t, util.ContainsIgnoreCase(slice, strings.ToLower(ele)))
		assert.True(t, util.ContainsIgnoreCase(slice, strings.ToUpper(ele)))
	}

	assert.False(t, util.Contains(slice, "not-existed"))

	assert.False(t, util.Contains([]string{}, "empty slice"))
}

func TestSliceUnique(t *testing.T) {
	duplicated := []string{"first", "second", "first", "third"}
	unique := util.Unique(duplicated)

	assert.Equal(t, 3, len(unique))
	assert.Equal(t, "first", unique[0])
	assert.Equal(t, "second", unique[1])
	assert.Equal(t, "third", unique[2])

	uniqueOfUnique := util.Unique(unique)

	assert.Equal(t, len(unique), len(uniqueOfUnique))
	for i, ele := range unique {
		assert.Equal(t, ele, uniqueOfUnique[i])
	}
}

func TestSliceRemove(t *testing.T) {
	slice := []string{"second", "first", "second", "third", "second", "second"}

	slice = util.Remove(slice, "second")

	assert.False(t, util.Contains(slice, "second"))
}
