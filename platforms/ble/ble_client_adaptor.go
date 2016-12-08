package ble

import (
	"fmt"
	"log"
	"strings"

	"github.com/currantlabs/gatt"
	"gobot.io/x/gobot"
)

var _ gobot.Adaptor = (*ClientAdaptor)(nil)

// Represents a Client Connection to a BLE Peripheral
type ClientAdaptor struct {
	name       string
	uuid       string
	device     gatt.Device
	peripheral gatt.Peripheral
	services   map[string]*Service
	connected  bool
	ready      chan struct{}
}

// NewClientAdaptor returns a new ClientAdaptor given a uuid
func NewClientAdaptor(uuid string) *ClientAdaptor {
	return &ClientAdaptor{
		name:      "BLECLient",
		uuid:      uuid,
		connected: false,
		ready:     make(chan struct{}),
		services:  make(map[string]*Service),
	}
}

func (b *ClientAdaptor) Name() string                { return b.name }
func (b *ClientAdaptor) SetName(n string)            { b.name = n }
func (b *ClientAdaptor) UUID() string                { return b.uuid }
func (b *ClientAdaptor) Peripheral() gatt.Peripheral { return b.peripheral }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (b *ClientAdaptor) Connect() (err error) {
	device, e := gatt.NewDevice(DefaultClientOptions...)
	if e != nil {
		log.Fatalf("Failed to open BLE device, err: %s\n", err)
		return e
	}

	b.device = device

	// Register handlers.
	device.Handle(
		gatt.PeripheralDiscovered(b.DiscoveryHandler),
		gatt.PeripheralConnected(b.ConnectHandler),
		gatt.PeripheralDisconnected(b.DisconnectHandler),
	)

	device.Init(b.StateChangeHandler)
	<-b.ready
	// TODO: make sure peripheral currently exists for this UUID before returning
	return nil
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
	b.peripheral.Device().CancelConnection(b.peripheral)

	return
}

// Finalize finalizes the BLEAdaptor
func (b *ClientAdaptor) Finalize() (err error) {
	return b.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) ReadCharacteristic(sUUID string, cUUID string) (data []byte, err error) {
	if !b.connected {
		log.Fatalf("Cannot read from BLE device until connected")
		return
	}

	characteristic := b.lookupCharacteristic(sUUID, cUUID)
	if characteristic == nil {
		log.Println("Cannot read from unknown characteristic")
		return
	}

	val, err := b.peripheral.ReadCharacteristic(characteristic)
	if err != nil {
		fmt.Printf("Failed to read characteristic, err: %s\n", err)
		return nil, err
	}

	return val, nil
}

// WriteCharacteristic writes bytes to the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) WriteCharacteristic(sUUID string, cUUID string, data []byte) (err error) {
	if !b.connected {
		log.Fatalf("Cannot write to BLE device until connected")
		return
	}

	characteristic := b.lookupCharacteristic(sUUID, cUUID)
	if characteristic == nil {
		log.Println("Cannot write to unknown characteristic")
		return
	}

	err = b.peripheral.WriteCharacteristic(characteristic, data, true)
	if err != nil {
		fmt.Printf("Failed to write characteristic, err: %s\n", err)
		return err
	}

	return
}

// Subscribe subscribes to notifications from the BLE device for the
// requested service and characteristic
func (b *ClientAdaptor) Subscribe(sUUID string, cUUID string, f func([]byte, error)) (err error) {
	if !b.connected {
		log.Fatalf("Cannot subscribe to BLE device until connected")
		return
	}

	characteristic := b.lookupCharacteristic(sUUID, cUUID)
	if characteristic == nil {
		log.Println("Cannot subscribe to unknown characteristic")
		return
	}

	fn := func(c *gatt.Characteristic, b []byte, err error) {
		f(b, err)
	}

	err = b.peripheral.SetNotifyValue(characteristic, fn)
	if err != nil {
		fmt.Printf("Failed to subscribe to characteristic, err: %s\n", err)
		return err
	}

	return
}

func (b *ClientAdaptor) StateChangeHandler(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func (b *ClientAdaptor) DiscoveryHandler(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	// try looking by local name
	if a.LocalName == b.UUID() {
		b.uuid = p.ID()
	} else {
		// try looking by ID
		id := strings.ToUpper(b.UUID())
		if strings.ToUpper(p.ID()) != id {
			return
		}
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	// and connect to it
	p.Device().Connect(p)
}

func (b *ClientAdaptor) ConnectHandler(p gatt.Peripheral, err error) {
	fmt.Printf("\nConnected Peripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())

	b.peripheral = p

	if err := p.SetMTU(250); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
	}

	ss, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

outer:
	for _, s := range ss {
		b.services[s.UUID().String()] = NewService(s.UUID().String(), s)

		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			continue
		}

		for _, c := range cs {
			_, err := p.DiscoverDescriptors(nil, c)
			if err != nil {
				fmt.Printf("Failed to discover descriptors: %v\n", err)
				continue outer
			}
			b.services[s.UUID().String()].characteristics[c.UUID().String()] = c
		}
	}

	b.connected = true
	close(b.ready)
}

func (b *ClientAdaptor) DisconnectHandler(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

// Finalize finalizes the ClientAdaptor
func (b *ClientAdaptor) lookupCharacteristic(sUUID string, cUUID string) *gatt.Characteristic {
	service := b.services[sUUID]
	if service == nil {
		log.Printf("Unknown service ID: %s\n", sUUID)
		return nil
	}

	characteristic := service.characteristics[cUUID]
	if characteristic == nil {
		log.Printf("Unknown characteristic ID: %s\n", cUUID)
		return nil
	}

	return characteristic
}

// Represents a BLE Peripheral's Service
type Service struct {
	uuid            string
	service         *gatt.Service
	characteristics map[string]*gatt.Characteristic
}

// NewService returns a new BLE Service given a uuid
func NewService(sUuid string, service *gatt.Service) *Service {
	return &Service{
		uuid:            sUuid,
		service:         service,
		characteristics: make(map[string]*gatt.Characteristic),
	}
}
