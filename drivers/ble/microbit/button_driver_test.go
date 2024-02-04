package microbit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

func TestNewButtonDriver(t *testing.T) {
	d := NewButtonDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &ButtonDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Button"))
	assert.NotNil(t, d.Eventer)
}

func TestButtonStartAndHalt(t *testing.T) {
	d := NewButtonDriver(testutil.NewBleTestAdaptor())
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestButtonReadData(t *testing.T) {
	sem := make(chan bool)
	a := testutil.NewBleTestAdaptor()
	d := NewButtonDriver(a)
	require.NoError(t, d.Start())

	err := d.On("buttonB", func(data interface{}) {
		sem <- true
	})
	require.NoError(t, err)

	a.SendTestDataToSubscriber([]byte{1}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"ButtonB\" was not published")
	}
}
