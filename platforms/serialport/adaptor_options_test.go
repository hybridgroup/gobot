package serialport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithName(t *testing.T) {
	// This is a general test, that options are applied by using the WithName() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange & act
	const newName = "new name"
	a := NewAdaptor("port", WithName(newName))
	// assert
	assert.Equal(t, newName, a.cfg.name)
}

func TestWithBaudRate(t *testing.T) {
	// arrange
	newBaudRate := 5432
	cfg := &configuration{baudRate: 1234}
	// act
	WithBaudRate(newBaudRate).apply(cfg)
	// assert
	assert.Equal(t, newBaudRate, cfg.baudRate)
}
