package gobot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommander(t *testing.T) {
	// arrange
	c := NewCommander()
	c.AddCommand("test", func(map[string]interface{}) interface{} {
		return "hi"
	})

	// act && assert
	assert.Equal(t, 1, len(c.Commands()))
	assert.NotNil(t, c.Command("test"))
	assert.Nil(t, c.Command("booyeah"))
}
