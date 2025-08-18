package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTracker(t *testing.T) {
	origin := os.Stdout
	tracker := NewFileTracker(os.Stdout)
	tracker.Begin()
	_, _ = fmt.Fprintln(os.Stdout, "hello world")
	assert.Equal(t, "hello world\n", tracker.Value())
	tracker.End()
	assert.Equal(t, "/dev/stdout", origin.Name())
	assert.Same(t, origin, os.Stdout)
}

func TestFileTrackerRef(t *testing.T) {
	origin := os.Stdout
	tracker := NewFileTracker(os.Stdout)
	tracker.Begin()
	_, _ = fmt.Fprintln(origin, "hello world")
	assert.Equal(t, "hello world\n", tracker.Value())
	tracker.End()
	assert.Equal(t, "/dev/stdout", origin.Name())
	assert.Same(t, origin, os.Stdout)
}
