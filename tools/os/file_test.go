package os

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExist(t *testing.T) {
	assert.True(t, IsDir(".."))
	assert.False(t, IsDir("none"))
	assert.True(t, Exist("file.go"))
}

func TestIsDir(t *testing.T) {
	assert.True(t, IsDir(".."))
	assert.False(t, IsDir("none"))
	assert.False(t, IsDir("file.go"))
}

func TestIsFile(t *testing.T) {
	assert.False(t, IsFile(".."))
	assert.False(t, IsFile("none"))
	assert.True(t, IsFile("file.go"))
}
