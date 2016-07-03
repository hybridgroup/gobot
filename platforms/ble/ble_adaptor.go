package ble

import (
	"fmt"
	"github.com/currantlabs/gatt"
	"github.com/hybridgroup/gobot"
	"log"
	"strings"
)

// TODO: handle other OS defaults besides Linux
var DefaultClientOptions = []gatt.Option{
	gatt.LnxMaxConnections(1),
	gatt.LnxDeviceID(-1, false),
}

var _ gobot.Adaptor = (*BLEAdaptor)(nil)

// Represents a Connection to a BLE Peripheral
type BLEAdaptor struct {
	name       string
	uuid       string
	device     gatt.Device
	peripheral gatt.Peripheral
	services   map[string]*BLEService
	connected  bool
	ready      chan struct{}
}

// NewBLEAdaptor returns a new BLEAdaptor given a name and uuid
func NewBLEAdaptor(name string, uuid string) *BLEAdaptor {
	return &BLEAdaptor{
		name:      name,
		uuid:      uuid,
		connected: false,
		ready:     make(chan struct{}),
		services:  make(map[string]*BLEService),
	}
}

func (b *BLEAdaptor) Name() string                { return b.name }
func (b *BLEAdaptor) UUID() string                { return b.uuid }
func (b *BLEAdaptor) Peripheral() gatt.Peripheral { return b.peripheral }

// Connect initiates a connection to the BLE peripheral. Returns true on successful connection.
func (b *BLEAdaptor) Connect() (errs []error) {
	device, err := gatt.NewDevice(DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open BLE device, err: %s\n", err)
		return
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
func (b *BLEAdaptor) Reconnect() (errs []error) {
	if b.connected {
		b.Disconnect()
	}
	return b.Connect()
}

// Disconnect terminates the connection to the BLE peripheral. Returns true on successful disconnect.
func (b *BLEAdaptor) Disconnect() (errs []error) {
	b.peripheral.Device().CancelConnection(b.peripheral)

	return
}

// Finalize finalizes the BLEAdaptor
func (b *BLEAdaptor) Finalize() (errs []error) {
	return b.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the
// requested service and characteristic
func (b *BLEAdaptor) ReadCharacteristic(sUUID string, cUUID string) (data []byte, err error) {
	if !b.connected {
		log.Fatalf("Cannot read from BLE device until connected")
		return
	}

	characteristic := b.services[sUUID].characteristics[cUUID]
	val, err := b.peripheral.ReadCharacteristic(characteristic)
	if err != nil {
		fmt.Printf("Failed to read characteristic, err: %s\n", err)
		return nil, err
	}

	return val, nil
}

func (b *BLEAdaptor) StateChangeHandler(d gatt.Device, s gatt.State) {
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

func (b *BLEAdaptor) DiscoveryHandler(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	id := strings.ToUpper(b.UUID())
	if strings.ToUpper(p.ID()) != id {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	// and connect to it
	p.Device().Connect(p)
}

func (b *BLEAdaptor) ConnectHandler(p gatt.Peripheral, err error) {
	fmt.Printf("\nConnected Peripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())

	b.peripheral = p

	if err := p.SetMTU(500); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
	}

	ss, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range ss {
		b.services[s.UUID().String()] = NewBLEService(s.UUID().String(), s)

		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			continue
		}

		for _, c := range cs {
			b.services[s.UUID().String()].characteristics[c.UUID().String()] = c
		}
	}

	b.connected = true
	close(b.ready)
	//defer p.Device().CancelConnection(p)
}

func (b *BLEAdaptor) DisconnectHandler(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

// Represents a BLE Peripheral's Service
type BLEService struct {
	uuid            string
	service         *gatt.Service
	characteristics map[string]*gatt.Characteristic
}

// NewBLEAdaptor returns a new BLEService given a uuid
func NewBLEService(sUuid string, service *gatt.Service) *BLEService {
	return &BLEService{
		uuid:            sUuid,
		service:         service,
		characteristics: make(map[string]*gatt.Characteristic),
	}
}
