package bebop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBebopValidatePitchWhenEqualOffset(t *testing.T) {
	assert.Equal(t, 100, ValidatePitch(32767.0, 32767.0))
}

func TestBebopValidatePitchWhenTiny(t *testing.T) {
	assert.Equal(t, 0, ValidatePitch(1.1, 32767.0))
}

func TestBebopValidatePitchWhenCentered(t *testing.T) {
	assert.Equal(t, 50, ValidatePitch(16383.5, 32767.0))
}
