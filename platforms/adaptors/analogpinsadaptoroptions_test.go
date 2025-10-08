package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func TestWithAnalogPinDebug(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val).apply(cfg)".
	// arrange & act
	a := NewAnalogPinsAdaptor(system.NewAccesser(), nil, WithAnalogPinDebug())
	// assert
	assert.True(t, a.analogPinsCfg.debug)
}
