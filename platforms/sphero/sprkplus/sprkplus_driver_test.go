package sprkplus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*SPRKPlusDriver)(nil)

func initTestSPRKPlusDriver() *SPRKPlusDriver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestSPRKPlusDriver(t *testing.T) {
	d := initTestSPRKPlusDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "SPRK"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestSPRKPlusDriverStartAndHalt(t *testing.T) {
	d := initTestSPRKPlusDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}
