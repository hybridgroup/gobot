package spi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriverWithStubbedAdaptor() (*Driver, *spiTestAdaptor) { //nolint:unparam // keep for further tests
	a := newSpiTestAdaptor()
	d := NewDriver(a, "SPI_BASIC")
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewDriver(t *testing.T) {
	var di interface{} = NewDriver(newSpiTestAdaptor(), "SPI_BASIC")
	d, ok := di.(*Driver)
	if !ok {
		t.Errorf("NewDriver() should have returned a *Driver")
	}
	assert.True(t, strings.HasPrefix(d.Name(), "SPI_BASIC"))
}

func TestStart(t *testing.T) {
	d := NewDriver(newSpiTestAdaptor(), "SPI_BASIC")
	assert.NoError(t, d.Start())
}

func TestHalt(t *testing.T) {
	d, _ := initTestDriverWithStubbedAdaptor()
	assert.NoError(t, d.Halt())
}

func TestSetName(t *testing.T) {
	// arrange
	d, _ := initTestDriverWithStubbedAdaptor()
	// act
	d.SetName("TESTME")
	// assert
	assert.Equal(t, "TESTME", d.Name())
}

func TestConnection(t *testing.T) {
	d, _ := initTestDriverWithStubbedAdaptor()
	assert.NotNil(t, d.Connection())
}
