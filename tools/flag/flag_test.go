package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DemoOptions struct {
	Help bool `name:"help" shorthand:"h" usage:"show help"`
}

func TestParse(t *testing.T) {
	options := &DemoOptions{}
	flagSet := ParseStruct(options, "demo", []string{"--help", "/other.exe", "-p"})
	assert.True(t, options.Help)
	assert.Equal(t, "demo", flagSet.Name())
	assert.EqualValues(t, []string{"/other.exe", "-p"}, flagSet.Args())
}
