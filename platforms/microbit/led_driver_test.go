package microbit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(NewBleTestAdaptor())
	return d
}

func TestLEDDriver(t *testing.T) {
	d := initTestLEDDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit LED"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestLEDDriverStartAndHalt(t *testing.T) {
	d := initTestLEDDriver()
	assert.Nil(t, d.Start())
	assert.Nil(t, d.Halt())
}

func TestLEDDriverWriteMatrix(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.Nil(t, d.WriteMatrix([]byte{0x01, 0x02}))
}

func TestLEDDriverWriteText(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.Nil(t, d.WriteText("Hello"))
}

func TestLEDDriverCommands(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.Nil(t, d.Blank())
	assert.Nil(t, d.Solid())
	assert.Nil(t, d.UpRightArrow())
	assert.Nil(t, d.UpLeftArrow())
	assert.Nil(t, d.DownRightArrow())
	assert.Nil(t, d.DownLeftArrow())
	assert.Nil(t, d.Dimond())
	assert.Nil(t, d.Smile())
	assert.Nil(t, d.Wink())
}
