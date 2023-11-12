package ardrone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArdroneValidatePitchWhenEqualOffset(t *testing.T) {
	assert.InDelta(t, 1.0, ValidatePitch(32767.0, 32767.0), 0.0)
}

func TestArdroneValidatePitchWhenTiny(t *testing.T) {
	assert.InDelta(t, 0.0, ValidatePitch(1.1, 32767.0), 0.0)
}

func TestArdroneValidatePitchWhenCentered(t *testing.T) {
	assert.InDelta(t, 0.5, ValidatePitch(16383.5, 32767.0), 0.0)
}
