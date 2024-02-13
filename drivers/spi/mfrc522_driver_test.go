package spi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MFRC522Driver)(nil)

func initTestMFRC522DriverWithStubbedAdaptor() (*MFRC522Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewMFRC522Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	// reset the written bytes during Start()
	a.spi.Reset()
	return d, a
}

func TestNewMFRC522Driver(t *testing.T) {
	var di interface{} = NewMFRC522Driver(newSpiTestAdaptor())
	d, ok := di.(*MFRC522Driver)
	if !ok {
		require.Fail(t, "NewMFRC522Driver() should have returned a *MFRC522Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "MFRC522"))
}

func TestMFRC522WriteByteData(t *testing.T) {
	// arrange
	d, a := initTestMFRC522DriverWithStubbedAdaptor()
	// act
	err := d.connection.WriteByteData(0x00, 0x00)
	// assert
	require.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00}, a.spi.Written())
}
