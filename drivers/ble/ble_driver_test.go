package ble

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver() *Driver {
	a := testutil.NewBleTestAdaptor()
	d := NewDriver(a, "BLE_BASIC", nil, nil)
	return d
}

func TestNewDriver(t *testing.T) {
	// arrange
	const name = "mybot"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewDriver(a, name, nil, nil)
	// assert
	assert.IsType(t, &Driver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), name))
	assert.Equal(t, a, d.Connection())
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func TestNewDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		name    = "mybot"
		newName = "overwrite mybot"
	)
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewDriver(a, name, nil, nil, WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func Test_applyWithName(t *testing.T) {
	// arrange
	const name = "mybot"
	cfg := configuration{name: "oldname"}
	// act
	WithName(name).apply(&cfg)
	// assert
	assert.Equal(t, name, cfg.name)
}

func TestStart(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Start())
	// arrange after start function
	d.afterStart = func() error { return fmt.Errorf("after start error") }
	// act, assert
	require.EqualError(t, d.Start(), "after start error")
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Halt())
	// arrange after start function
	d.beforeHalt = func() error { return fmt.Errorf("before halt error") }
	// act, assert
	require.EqualError(t, d.Halt(), "before halt error")
}

func TestAdaptor(t *testing.T) {
	wrongConnectorType := struct {
		a uint32
	}{}
	// arrange
	a := testutil.NewBleTestAdaptor()
	d := NewDriver(a, "BLE_BASIC", nil, nil)
	// act, assert
	assert.Equal(t, a, d.Adaptor())
	// arrange wrong connector type
	d.connection = wrongConnectorType
	// act, assert
	assert.Nil(t, d.Adaptor())
}
