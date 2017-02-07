package ardrone

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestArdroneValidatePitchWhenEqualOffset(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(32767.0, 32767.0), 1.0)
}

func TestArdroneValidatePitchWhenTiny(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(1.1, 32767.0), 0.0)
}

func TestArdroneValidatePitchWhenCentered(t *testing.T) {
	gobottest.Assert(t, ValidatePitch(16383.5, 32767.0), 0.5)
}
