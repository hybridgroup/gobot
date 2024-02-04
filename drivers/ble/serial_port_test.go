package ble

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*SerialPortDriver)(nil)

var _ io.ReadWriteCloser = (*SerialPortDriver)(nil)

func TestBLESerialPort(t *testing.T) {
	d := NewSerialPortDriver(testutil.NewBleTestAdaptor(), "123", "456")
	assert.Equal(t, "01:02:03:0A:0B:0C", d.Address())
}
