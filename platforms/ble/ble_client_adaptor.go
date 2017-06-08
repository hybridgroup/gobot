package ble

import (
	"context"
	"log"
	"strings"
	"sync"

	"gobot.io/x/gobot"

	blelib "github.com/currantlabs/ble"
	"github.com/pkg/errors"
)

var currentDevice *blelib.Device
var bleMutex sync.Mutex
var bleCtx context.Context

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
	WithoutReponses(bool)
}

// ClientAdaptor represents a Client Connection to a BLE Peripheral
type ClientAdaptor struct {
	name       string
	address    string
	DeviceName string

	addr    blelib.Addr
	device  *blelib.Device
	client  blelib.Client
	profile *blelib.Profile

	connected       bool
	ready           chan struct{}
	withoutReponses bool
}

// NewClientAdaptor returns a new ClientAdaptor given an address or peripheral name
func NewClientAdaptor(address string) *ClientAdaptor {
	return &ClientAdaptor{
		name:            gobot.DefaultName("BLEClient"),
		address:         address,
		DeviceName:      "default",
		connected:       false,
		withoutReponses: false,
	}
}

// Name returns the name for the adaptor
func (b *ClientAdaptor) Name() string { return b.name }

// SetName sets the name for the adaptor
func (b *ClientAdaptor) SetName(n string) { b.name = n }

// Address returns the Bluetooth LE address for the adaptor
func (b *ClientAdaptor) Address() string { return b.address }

// Address sets if the adaptor should expect responses after
// writing characteristics for this device
func (b *ClientAdaptor) WithoutReponses(use bool) { b.withoutReponses = use }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (b *ClientAdaptor) Connect() (err error) {
	bleMutex.Lock()
	defer bleMutex.Unlock()

	b.device, err = getBLEDevice(b.DeviceName)
	if err != nil {
		return errors.Wrap(err, "can't connect to device "+b.DeviceName)
	}

	var cln blelib.Client

	cln, err = blelib.Connect(context.Background(), filter(b.Address()))
	if err != nil {
		return errors.Wrap(err, "can't connect to peripheral "+b.Address())
	}

	b.addr = cln.Address()
	b.address = cln.Address().String()
	b.SetName(cln.Name())
	b.client = cln

	p, err := b.client.DiscoverProfile(true)
	if err != nil {
		return errors.Wrap(err, "can't discover profile")
	}

	b.profile = p
	b.connected = true
	return
}

// Reconnect attempts to reconnect to the BLE peripheral. If it has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (b *ClientAdaptor) Reconnect() (err error) {
	if b.connected {
		b.Disconnect()
	}
	return b.Connect()
}

// Disconnect terminates the connection to the BLE peripheral. Returns true on successful disconnect.
func (b *ClientAdaptor) Disconnect() (err error) {
	b.client.CancelConnection()
	return
}

// Finalize finalizes the BLEAdaptor
func (b *ClientAdaptor) Finalize() (err error) {
	return b.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the
// requested characteristic uuid
func (b *ClientAdaptor) ReadCharacteristic(cUUID string) (data []byte, err error) {
	if !b.connected {
		log.Fatalf("Cannot read from BLE device until connected")
		return
	}

	uuid, _ := blelib.Parse(cUUID)

	if u := b.profile.Find(blelib.NewCharacteristic(uuid)); u != nil {
		data, err = b.client.ReadCharacteristic(u.(*blelib.Characteristic))
	}

	return
}

// WriteCharacteristic writes bytes to the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) WriteCharacteristic(cUUID string, data []byte) (err error) {
	if !b.connected {
		log.Println("Cannot write to BLE device until connected")
		return
	}

	uuid, _ := blelib.Parse(cUUID)

	if u := b.profile.Find(blelib.NewCharacteristic(uuid)); u != nil {
		err = b.client.WriteCharacteristic(u.(*blelib.Characteristic), data, b.withoutReponses)
	}

	return
}

// Subscribe subscribes to notifications from the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) (err error) {
	if !b.connected {
		log.Fatalf("Cannot subscribe to BLE device until connected")
		return
	}

	uuid, _ := blelib.Parse(cUUID)

	if u := b.profile.Find(blelib.NewCharacteristic(uuid)); u != nil {
		h := func(req []byte) { f(req, nil) }
		err = b.client.Subscribe(u.(*blelib.Characteristic), false, h)
		if err != nil {
			return err
		}
		return nil
	}

	return
}

// getBLEDevice is singleton for blelib HCI device connection
func getBLEDevice(impl string) (d *blelib.Device, err error) {
	if currentDevice != nil {
		return currentDevice, nil
	}

	dev, e := defaultDevice(impl)
	if e != nil {
		return nil, errors.Wrap(e, "can't get device")
	}
	blelib.SetDefaultDevice(dev)

	currentDevice = &dev
	d = &dev
	return
}

func filter(name string) blelib.AdvFilter {
	return func(a blelib.Advertisement) bool {
		return strings.ToLower(a.LocalName()) == strings.ToLower(name) ||
			a.Address().String() == strings.ToLower(name)
	}
}
