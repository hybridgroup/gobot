package bleclient

import (
	"fmt"
	"log"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"

	"gobot.io/x/gobot/v2"
)

type configuration struct {
	scanTimeout          time.Duration
	sleepAfterDisconnect time.Duration
	debug                bool
}

// Adaptor represents a Client Connection to a BLE Peripheral
type Adaptor struct {
	name       string
	identifier string
	cfg        *configuration

	btAdpt          *btAdapter
	btDevice        *btDevice
	characteristics map[string]bluetoothExtCharacteristicer

	connected bool
	rssi      int

	btAdptCreator btAdptCreatorFunc
	mutex         *sync.Mutex
}

// NewAdaptor returns a new Adaptor given an identifier. The identifier can be the address or the name.
//
// Supported options:
//
//	"WithAdaptorDebug"
//	"WithAdaptorScanTimeout"
func NewAdaptor(identifier string, opts ...optionApplier) *Adaptor {
	cfg := configuration{
		scanTimeout:          10 * time.Minute,
		sleepAfterDisconnect: 500 * time.Millisecond,
	}

	a := Adaptor{
		name:            gobot.DefaultName("BLEClient"),
		identifier:      identifier,
		cfg:             &cfg,
		characteristics: make(map[string]bluetoothExtCharacteristicer),
		btAdptCreator:   newBtAdapter,
		mutex:           &sync.Mutex{},
	}

	for _, o := range opts {
		o.apply(a.cfg)
	}

	return &a
}

// WithDebug switch on some debug messages.
func WithDebug() debugOption {
	return debugOption(true)
}

// WithScanTimeout substitute the default scan timeout of 10 min.
func WithScanTimeout(timeout time.Duration) scanTimeoutOption {
	return scanTimeoutOption(timeout)
}

// Name returns the name for the adaptor and after the connection is done, the name of the device
func (a *Adaptor) Name() string {
	if a.btDevice != nil {
		return a.btDevice.name()
	}
	return a.name
}

// SetName sets the name for the adaptor
func (a *Adaptor) SetName(n string) { a.name = n }

// Address returns the Bluetooth LE address of the device if connected, otherwise the identifier
func (a *Adaptor) Address() string {
	if a.btDevice != nil {
		return a.btDevice.address()
	}

	return a.identifier
}

// RSSI returns the Bluetooth LE RSSI value at the moment of connecting the adaptor
func (a *Adaptor) RSSI() int { return a.rssi }

// WithoutResponses sets if the adaptor should expect responses after
// writing characteristics for this device (has no effect at the moment).
func (a *Adaptor) WithoutResponses(bool) {}

// Connect initiates a connection to the BLE peripheral.
func (a *Adaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error

	if a.cfg.debug {
		fmt.Println("[Connect]: enable adaptor...")
	}

	// for re-connect, the adapter is already known
	if a.btAdpt == nil {
		a.btAdpt = a.btAdptCreator(bluetooth.DefaultAdapter, a.cfg.debug)
		if err := a.btAdpt.enable(); err != nil {
			return fmt.Errorf("can't get adapter default: %w", err)
		}
	}

	if a.cfg.debug {
		fmt.Printf("[Connect]: scan %s for the identifier '%s'...\n", a.cfg.scanTimeout, a.identifier)
	}

	result, err := a.btAdpt.scan(a.identifier, a.cfg.scanTimeout)
	if err != nil {
		return err
	}

	if a.cfg.debug {
		fmt.Printf("[Connect]: connect to peripheral device with address %s...\n", result.Address)
	}

	dev, err := a.btAdpt.connect(result.Address, result.LocalName())
	if err != nil {
		return err
	}

	a.rssi = int(result.RSSI)
	a.btDevice = dev

	if a.cfg.debug {
		fmt.Println("[Connect]: get all services/characteristics...")
	}
	services, err := a.btDevice.discoverServices(nil)
	if err != nil {
		return err
	}
	for _, service := range services {
		if a.cfg.debug {
			fmt.Printf("[Connect]: service found: %s\n", service)
		}
		chars, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, char := range chars {
			if a.cfg.debug {
				fmt.Printf("[Connect]: characteristic found: %s\n", char)
			}
			c := char // to prevent implicit memory aliasing in for loop, before go 1.22
			a.characteristics[char.UUID().String()] = &c
		}
	}

	if a.cfg.debug {
		fmt.Println("[Connect]: connected")
	}
	a.connected = true
	return nil
}

// Reconnect attempts to reconnect to the BLE peripheral. If it has an active connection
// it will first close that connection and then establish a new connection.
func (a *Adaptor) Reconnect() error {
	if a.connected {
		if err := a.Disconnect(); err != nil {
			return err
		}
	}
	return a.Connect()
}

// Disconnect terminates the connection to the BLE peripheral.
func (a *Adaptor) Disconnect() error {
	if a.cfg.debug {
		fmt.Println("[Disconnect]: disconnect...")
	}
	err := a.btDevice.disconnect()
	time.Sleep(a.cfg.sleepAfterDisconnect)
	a.connected = false
	if a.cfg.debug {
		fmt.Println("[Disconnect]: disconnected")
	}
	return err
}

// Finalize finalizes the BLEAdaptor
func (a *Adaptor) Finalize() error {
	return a.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the requested characteristic UUID.
// The UUID can be given as 16-bit or 128-bit (with or without dashes) value.
func (a *Adaptor) ReadCharacteristic(cUUID string) ([]byte, error) {
	if !a.connected {
		return nil, fmt.Errorf("cannot read from BLE device until connected")
	}

	cUUID, err := convertUUID(cUUID)
	if err != nil {
		return nil, err
	}

	if chara, ok := a.characteristics[cUUID]; ok {
		return readFromCharacteristic(chara)
	}

	return nil, fmt.Errorf("unknown characteristic: %s", cUUID)
}

// WriteCharacteristic writes bytes to the BLE device for the requested characteristic UUID.
// The UUID can be given as 16-bit or 128-bit (with or without dashes) value.
func (a *Adaptor) WriteCharacteristic(cUUID string, data []byte) error {
	if !a.connected {
		return fmt.Errorf("cannot write to BLE device until connected")
	}

	cUUID, err := convertUUID(cUUID)
	if err != nil {
		return err
	}

	if chara, ok := a.characteristics[cUUID]; ok {
		return writeToCharacteristicWithoutResponse(chara, data)
	}

	return fmt.Errorf("unknown characteristic: %s", cUUID)
}

// Subscribe subscribes to notifications from the BLE device for the requested characteristic UUID.
// The UUID can be given as 16-bit or 128-bit (with or without dashes) value.
func (a *Adaptor) Subscribe(cUUID string, f func(data []byte)) error {
	if !a.connected {
		return fmt.Errorf("cannot subscribe to BLE device until connected")
	}

	cUUID, err := convertUUID(cUUID)
	if err != nil {
		return err
	}

	if chara, ok := a.characteristics[cUUID]; ok {
		return enableNotificationsForCharacteristic(chara, f)
	}

	return fmt.Errorf("unknown characteristic: %s", cUUID)
}
