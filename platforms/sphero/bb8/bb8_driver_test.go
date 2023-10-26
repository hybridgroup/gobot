package bb8

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*BB8Driver)(nil)

func initTestBB8Driver() *BB8Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestBB8Driver(t *testing.T) {
	d := initTestBB8Driver()
	assert.True(t, strings.HasPrefix(d.Name(), "BB8"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestBB8DriverStartAndHalt(t *testing.T) {
	d := initTestBB8Driver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}
