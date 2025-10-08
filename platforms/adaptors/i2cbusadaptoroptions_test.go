package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func TestWithI2cDebug(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val).apply(cfg)".
	// arrange & act
	a := NewI2cBusAdaptor(system.NewAccesser(), nil, 0, WithI2cDebug())
	// assert
	assert.True(t, a.i2cBusCfg.debug)
}
