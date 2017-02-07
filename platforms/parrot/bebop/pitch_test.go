package bebop

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestBebopValidatePitchWhenEqualOffset(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(32767.0, 32767.0), 100)
}

func TestBebopValidatePitchWhenTiny(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(1.1, 32767.0), 0)
}

func TestBebopValidatePitchWhenCentered(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(16383.5, 32767.0), 50)
}
