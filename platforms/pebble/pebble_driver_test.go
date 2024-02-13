package pebble

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver() *Driver {
	return NewDriver(NewAdaptor())
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	require.NoError(t, d.Start())
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	require.NoError(t, d.Halt())
}

func TestDriver(t *testing.T) {
	d := initTestDriver()

	assert.Equal(t, "Pebble", d.Name())
	assert.Equal(t, "Pebble", d.Connection().Name())

	sem := make(chan bool)
	d.SendNotification("Hello")
	d.SendNotification("World")

	assert.Equal(t, "Hello", d.Messages[0])
	assert.Equal(t, "Hello", d.PendingMessage())
	assert.Equal(t, "World", d.PendingMessage())
	assert.Equal(t, "", d.PendingMessage())

	_ = d.On(d.Event("button"), func(data interface{}) {
		sem <- true
	})

	d.PublishEvent("button", "")

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Button Event was not published")
	}

	_ = d.On(d.Event("accel"), func(data interface{}) {
		sem <- true
	})

	d.Command("publish_event")(map[string]interface{}{"name": "accel", "data": "100"})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Accel Event was not published")
	}

	d.Command("send_notification")(map[string]interface{}{"message": "Hey buddy!"})
	assert.Equal(t, "Hey buddy!", d.Messages[0])

	message := d.Command("pending_message")(map[string]interface{}{})
	assert.Equal(t, "Hey buddy!", message)
}
