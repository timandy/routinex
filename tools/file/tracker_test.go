package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTracker(t *testing.T) {
	origin := os.Stdout
	tracker := NewFileTracker(&os.Stdout)
	tracker.Begin()
	fmt.Println("hello world")
	assert.Equal(t, "hello world\n", tracker.Value())
	tracker.End()
	assert.Same(t, origin, os.Stdout)
}
