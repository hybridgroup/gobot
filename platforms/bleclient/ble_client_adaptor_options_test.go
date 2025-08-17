package bleclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDebug(t *testing.T) {
	// This is a general test, that options are applied by using the WithDebug() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange & act
	a := NewAdaptor("address", WithDebug())
	// assert
	assert.True(t, a.cfg.debug)
}

func TestWithDropCharacteristicsOnDisconnect(t *testing.T) {
	// arrange
	cfg := &configuration{dropCharacteristicsOnDisconnect: false}
	// act
	WithDropCharacteristicsOnDisconnect().apply(cfg)
	// assert
	assert.True(t, cfg.dropCharacteristicsOnDisconnect)
}

func TestWithScanTimeout(t *testing.T) {
	// arrange
	newTimeout := 2 * time.Second
	cfg := &configuration{scanTimeout: 10 * time.Second}
	// act
	WithScanTimeout(newTimeout).apply(cfg)
	// assert
	assert.Equal(t, newTimeout, cfg.scanTimeout)
}
