package system

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newSpiGpio(t *testing.T) {
	// arrange
	dpa := newMockDigitalPinAccess(nil)
	cfg := spiGpioConfig{
		pinProvider: dpa,
		sclkPinID:   "1",
		ncsPinID:    "2",
		sdoPinID:    "3",
		sdiPinID:    "4",
	}
	// act
	d, err := newSpiGpio(cfg, 10001)
	// assert
	require.NoError(t, err)
	assert.Equal(t, cfg, d.cfg)
	assert.Equal(t, 50*time.Microsecond, d.tclk)
	assert.NotNil(t, d.sclkPin)
	assert.NotNil(t, d.ncsPin)
	assert.NotNil(t, d.sdoPin)
	assert.NotNil(t, d.sdiPin)
	assert.Equal(t, 1, dpa.AppliedOptions("", cfg.sclkPinID))
	assert.Equal(t, 1, dpa.AppliedOptions("", cfg.ncsPinID))
	assert.Equal(t, 1, dpa.AppliedOptions("", cfg.sdoPinID))
	assert.Equal(t, 0, dpa.AppliedOptions("", cfg.sdiPinID)) // already input, so no option applied
}

func TestSpiGpioTxRx(t *testing.T) {
	// arrange
	dpa := newMockDigitalPinAccess(nil)
	cfg := spiGpioConfig{
		pinProvider: dpa,
		sclkPinID:   "1",
		ncsPinID:    "2",
		sdoPinID:    "3",
		sdiPinID:    "4",
	}
	dpa.UseValues("", cfg.sdiPinID, []int{1, 0, 0, 1, 0, 1, 0, 1})
	d, err := newSpiGpio(cfg, 10001)
	require.NoError(t, err)
	// act
	rx := []byte{0x87}
	err = d.TxRx([]byte{0xf1}, rx)
	// assert
	require.NoError(t, err)
	assert.Equal(t, []int{0, 1}, dpa.Written("", cfg.ncsPinID))
	assert.Equal(t, []int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, dpa.Written("", cfg.sclkPinID))
	assert.Equal(t, []int{128, 64, 32, 16, 0, 0, 0, 1}, dpa.Written("", cfg.sdoPinID)) // > 0 just means a 1 was send
	assert.Empty(t, dpa.Written("", cfg.sdiPinID))
	assert.Equal(t, []byte{0x95}, rx)
}

func TestSpiGpioClose(t *testing.T) {
	// arrange
	dpa := newMockDigitalPinAccess(nil)
	cfg := spiGpioConfig{
		pinProvider: dpa,
		sclkPinID:   "1",
		ncsPinID:    "2",
		sdoPinID:    "3",
		sdiPinID:    "4",
	}
	dpa.UseValues("", cfg.sdiPinID, []int{1, 0, 0, 1, 0, 1, 0, 1})
	d, err := newSpiGpio(cfg, 10001)
	require.NoError(t, err)
	// act
	err = d.Close()
	// assert
	require.NoError(t, err)
	assert.Equal(t, -1, dpa.Exported("", cfg.sclkPinID))
	assert.Equal(t, -1, dpa.Exported("", cfg.ncsPinID))
	assert.Equal(t, -1, dpa.Exported("", cfg.sdoPinID))
	assert.Equal(t, -1, dpa.Exported("", cfg.sdiPinID))
}
