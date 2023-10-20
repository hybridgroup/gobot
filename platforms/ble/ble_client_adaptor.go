package ble

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"

	"tinygo.org/x/bluetooth"
)

var (
	currentAdapter *bluetooth.Adapter
	bleMutex       sync.Mutex
)

// BLEConnector is the interface that a BLE ClientAdaptor must implement
type BLEConnector interface {
	Connect() error
	Reconnect() error
	Disconnect() error
	Finalize() error
	Name() string
	SetName(string)
	Address() string
	ReadCharacteristic(string) ([]byte, error)
	WriteCharacteristic(string, []byte) error
	Subscribe(string, func([]byte, error)) error
	WithoutResponses(bool)
}

// ClientAdaptor represents a Client Connection to a BLE Peripheral
type ClientAdaptor struct {
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

// NewClientAdaptor returns a new ClientAdaptor given an address
func NewClientAdaptor(address string) *ClientAdaptor {
	return &ClientAdaptor{
		name:             gobot.DefaultName("BLEClient"),
		address:          address,
		AdapterName:      "default",
		connected:        false,
		withoutResponses: false,
		characteristics:  make(map[string]bluetooth.DeviceCharacteristic),
	}
}

// Name returns the name for the adaptor
func (b *ClientAdaptor) Name() string { return b.name }

// SetName sets the name for the adaptor
func (b *ClientAdaptor) SetName(n string) { b.name = n }

// Address returns the Bluetooth LE address for the adaptor
func (b *ClientAdaptor) Address() string { return b.address }

// WithoutResponses sets if the adaptor should expect responses after
// writing characteristics for this device
func (b *ClientAdaptor) WithoutResponses(use bool) { b.withoutResponses = use }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (b *ClientAdaptor) Connect() error {
	bleMutex.Lock()
	defer bleMutex.Unlock()

	var err error
	// enable adaptor
	b.adpt, err = getBLEAdapter(b.AdapterName)
	if err != nil {
		return fmt.Errorf("can't get adapter %s: %w", b.AdapterName, err)
	}

	// handle address
	b.addr.Set(b.Address())

	// scan for the address
	ch := make(chan bluetooth.ScanResult, 1)
	err = b.adpt.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == b.Address() {
			if err := b.adpt.StopScan(); err != nil {
				panic(err)
			}
			b.SetName(result.LocalName())
			ch <- result
		}
	})

	if err != nil {
		return err
	}

	// wait to connect to peripheral device
	result := <-ch
	b.device, err = b.adpt.Connect(result.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	// get all services/characteristics
	srvcs, err := b.device.DiscoverServices(nil)
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
			b.characteristics[char.UUID().String()] = char
		}
	}

	b.connected = true
	return nil
}

// Reconnect attempts to reconnect to the BLE peripheral. If it has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (b *ClientAdaptor) Reconnect() error {
	if b.connected {
		if err := b.Disconnect(); err != nil {
			return err
		}
	}
	return b.Connect()
}

// Disconnect terminates the connection to the BLE peripheral. Returns true on successful disconnect.
func (b *ClientAdaptor) Disconnect() error {
	err := b.device.Disconnect()
	time.Sleep(500 * time.Millisecond)
	return err
}

// Finalize finalizes the BLEAdaptor
func (b *ClientAdaptor) Finalize() error {
	return b.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the
// requested characteristic uuid
func (b *ClientAdaptor) ReadCharacteristic(cUUID string) (data []byte, err error) {
	if !b.connected {
		return nil, fmt.Errorf("Cannot read from BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := b.characteristics[cUUID]; ok {
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
func (b *ClientAdaptor) WriteCharacteristic(cUUID string, data []byte) error {
	if !b.connected {
		return fmt.Errorf("Cannot write to BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := b.characteristics[cUUID]; ok {
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
func (b *ClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) error {
	if !b.connected {
		return fmt.Errorf("Cannot subscribe to BLE device until connected")
	}

	cUUID = convertUUID(cUUID)

	if char, ok := b.characteristics[cUUID]; ok {
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
