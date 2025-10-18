package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func TestWithOneWireDebug(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val).apply(cfg)".
	// arrange & act
	a := NewOneWireBusAdaptor(system.NewAccesser(), WithOneWireDebug())
	// assert
	assert.True(t, a.oneWireBusCfg.debug)
}
