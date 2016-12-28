package ble

import (
	"context"
	"fmt"
	"log"
	"strings"

	blelib "github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
	"github.com/pkg/errors"
)

// ClientAdaptor represents a Client Connection to a BLE Peripheral
type ClientAdaptor struct {
	name    string
	address string

	//	uuid    blelib.UUID
	addr    blelib.Addr
	device  blelib.Device
	client  blelib.Client
	profile *blelib.Profile

	connected bool
	ready     chan struct{}
}

// NewClientAdaptor returns a new ClientAdaptor given an address or peripheral name
func NewClientAdaptor(address string) *ClientAdaptor {
	return &ClientAdaptor{
		name:      "BLECLient",
		address:   address,
		connected: false,
	}
}

func (b *ClientAdaptor) Name() string     { return b.name }
func (b *ClientAdaptor) SetName(n string) { b.name = n }
func (b *ClientAdaptor) Address() string  { return b.address }

//func (b *ClientAdaptor) Peripheral() gatt.Peripheral { return b.peripheral }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (b *ClientAdaptor) Connect() (err error) {
	d, err := dev.NewDevice("default")
	if err != nil {
		return errors.Wrap(err, "can't new device")
	}
	blelib.SetDefaultDevice(d)
	b.device = d

	var cln blelib.Client

	ctx := blelib.WithSigHandler(context.WithTimeout(context.Background(), 0))
	cln, err = blelib.Connect(ctx, filter(b.Name()))
	if err != nil {
		return errors.Wrap(err, "can't connect")
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
		if err != nil {
			return nil, err
		}
		fmt.Printf("    Value         %x | %q\n", data, data)
		return data, nil
	}

	return data, nil
}

// WriteCharacteristic writes bytes to the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) WriteCharacteristic(cUUID string, data []byte) (err error) {
	if !b.connected {
		log.Fatalf("Cannot write to BLE device until connected")
		return
	}

	uuid, _ := blelib.Parse(cUUID)

	if u := b.profile.Find(blelib.NewCharacteristic(uuid)); u != nil {
		err = b.client.WriteCharacteristic(u.(*blelib.Characteristic), data, false)
		if err != nil {
			return err
		}
		return nil
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

func filter(name string) blelib.AdvFilter {
	return func(a blelib.Advertisement) bool {
		return strings.ToLower(a.LocalName()) == strings.ToLower(name) ||
			a.Address().String() == strings.ToLower(name)
	}
}
