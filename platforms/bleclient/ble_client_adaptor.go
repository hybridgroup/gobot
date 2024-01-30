package bleclient

import (
	"fmt"
	"log"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"

	"gobot.io/x/gobot/v2"
)

var (
	currentAdapter *bluetooth.Adapter
	bleMutex       sync.Mutex
)

// Adaptor represents a client connection to a BLE Peripheral
type Adaptor struct {
	name        string
	address     string
	AdapterName string

	addr            bluetooth.Address
	adpt            *bluetooth.Adapter
	device          *bluetooth.Device
	characteristics map[string]bluetooth.DeviceCharacteristic

	connected        bool
	withoutResponses bool
}

// NewAdaptor returns a new Bluetooth LE client adaptor given an address
func NewAdaptor(address string) *Adaptor {
	return &Adaptor{
		name:             gobot.DefaultName("BLEClient"),
		address:          address,
		AdapterName:      "default",
		connected:        false,
		withoutResponses: false,
		characteristics:  make(map[string]bluetooth.DeviceCharacteristic),
	}
}

// Name returns the name for the adaptor
func (a *Adaptor) Name() string { return a.name }

// SetName sets the name for the adaptor
func (a *Adaptor) SetName(n string) { a.name = n }

// Address returns the Bluetooth LE address for the adaptor
func (a *Adaptor) Address() string { return a.address }

// WithoutResponses sets if the adaptor should expect responses after
// writing characteristics for this device
func (a *Adaptor) WithoutResponses(use bool) { a.withoutResponses = use }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (a *Adaptor) Connect() error {
	bleMutex.Lock()
	defer bleMutex.Unlock()

	var err error
	// enable adaptor
	a.adpt, err = getBLEAdapter(a.AdapterName)
	if err != nil {
		return fmt.Errorf("can't get adapter %s: %w", a.AdapterName, err)
	}

	// handle address
	a.addr.Set(a.Address())

	// scan for the address
	ch := make(chan bluetooth.ScanResult, 1)
	err = a.adpt.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == a.Address() {
			if err := a.adpt.StopScan(); err != nil {
				panic(err)
			}
			a.SetName(result.LocalName())
			ch <- result
		}
	})

	if err != nil {
		return err
	}

	// wait to connect to peripheral device
	result := <-ch
	a.device, err = a.adpt.Connect(result.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	// get all services/characteristics
	srvcs, err := a.device.DiscoverServices(nil)
	if err != nil {
		return err
	}
	for _, srvc := range srvcs {
		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, char := range chars {
			a.characteristics[char.UUID().String()] = char
		}
	}

	a.connected = true
	return nil
}

// Reconnect attempts to reconnect to the BLE peripheral. If it has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (a *Adaptor) Reconnect() error {
	if a.connected {
		if err := a.Disconnect(); err != nil {
			return err
		}
	}
	return a.Connect()
}

// Disconnect terminates the connection to the BLE peripheral. Returns true on successful disconnect.
func (a *Adaptor) Disconnect() error {
	err := a.device.Disconnect()
	time.Sleep(500 * time.Millisecond)
	return err
}

// Finalize finalizes the BLEAdaptor
func (a *Adaptor) Finalize() error {
	return a.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the
// requested characteristic uuid
func (a *Adaptor) ReadCharacteristic(cUUID string) ([]byte, error) {
	if !a.connected {
		return nil, fmt.Errorf("Cannot read from BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := a.characteristics[cUUID]; ok {
		buf := make([]byte, 255)
		n, err := char.Read(buf)
		if err != nil {
			return nil, err
		}
		return buf[:n], nil
	}

	return nil, fmt.Errorf("Unknown characteristic: %s", cUUID)
}

// WriteCharacteristic writes bytes to the BLE device for the
// requested service and characteristic
func (a *Adaptor) WriteCharacteristic(cUUID string, data []byte) error {
	if !a.connected {
		return fmt.Errorf("Cannot write to BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := a.characteristics[cUUID]; ok {
		_, err := char.WriteWithoutResponse(data)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Unknown characteristic: %s", cUUID)
}

// Subscribe subscribes to notifications from the BLE device for the
// requested service and characteristic
func (a *Adaptor) Subscribe(cUUID string, f func([]byte, error)) error {
	if !a.connected {
		return fmt.Errorf("Cannot subscribe to BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := a.characteristics[cUUID]; ok {
		fn := func(d []byte) {
			f(d, nil)
		}
		return char.EnableNotifications(fn)
	}

	return fmt.Errorf("Unknown characteristic: %s", cUUID)
}

// getBLEAdapter is singleton for bluetooth adapter connection
func getBLEAdapter(impl string) (*bluetooth.Adapter, error) { //nolint:unparam // TODO: impl is unused, maybe an error
	if currentAdapter != nil {
		return currentAdapter, nil
	}

	currentAdapter = bluetooth.DefaultAdapter
	err := currentAdapter.Enable()
	if err != nil {
		return nil, err
	}

	return currentAdapter, nil
}
