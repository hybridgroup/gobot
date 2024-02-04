//go:build !windows
// +build !windows

package firmata

import (
	"io"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

const (
	// ReceiveID is the BLE characteristic ID for receiving serial data
	ReceiveID = "6e400003b5a3f393e0a9e50e24dcca9e"

	// TransmitID is the BLE characteristic ID for transmitting serial data
	TransmitID = "6e400002b5a3f393e0a9e50e24dcca9e"
)

// BLEAdaptor represents a Bluetooth LE based connection to a
// microcontroller running FirmataBLE
type BLEAdaptor struct {
	*Adaptor
}

// NewBLEAdaptor opens and uses a BLE connection to a
// microcontroller running FirmataBLE
func NewBLEAdaptor(args ...interface{}) *BLEAdaptor {
	address := args[0].(string) //nolint:forcetypeassert // ok here
	rid := ReceiveID
	wid := TransmitID

	//nolint:forcetypeassert // ok here
	if len(args) >= 3 {
		rid = args[1].(string)
		wid = args[2].(string)
	}

	a := NewAdaptor(address)
	a.SetName(gobot.DefaultName("BLEFirmata"))
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		a := bleclient.NewAdaptor(address)
		sp := ble.NewSerialPortDriver(a, rid, wid)
		if err := sp.Open(); err != nil {
			return sp, err
		}
		return sp, nil
	}

	return &BLEAdaptor{
		Adaptor: a,
	}
}
