package ble

import (
	"fmt"
	"log"
	"strings"
	"github.com/hybridgroup/gobot"
	"github.com/paypal/gatt"
)

// TODO: handle other OS defaults besides Linux
var DefaultClientOptions = []gatt.Option{
	gatt.LnxMaxConnections(1),
	gatt.LnxDeviceID(-1, false),
}

var _ gobot.Adaptor = (*BLEAdaptor)(nil)

// Represents a Connection to a BLE Peripheral
type BLEAdaptor struct {
	name      string
	uuid      string
	device    gatt.Device
	peripheral 				gatt.Peripheral
	//sp        io.ReadWriteCloser
	connected bool
	//connect   func(string) (io.ReadWriteCloser, error)
}

// NewBLEAdaptor returns a new BLEAdaptor given a name and uuid
func NewBLEAdaptor(name string, uuid string) *BLEAdaptor {
	return &BLEAdaptor{
		name: name,
		uuid: uuid,
		connected: false,
		// connect: func(port string) (io.ReadWriteCloser, error) {
		// 	return serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
		// },
	}
}

func (b *BLEAdaptor) Name() string { return b.name }
func (b *BLEAdaptor) UUID() string { return b.uuid }

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
		gatt.PeripheralDiscovered(b.onDiscovered),
		gatt.PeripheralConnected(b.onConnected),
		gatt.PeripheralDisconnected(b.onDisconnected),
	)

	device.Init(b.onStateChanged)

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
	// if a.connected {
	// 	if err := a.sp.Close(); err != nil {
	// 		return []error{err}
	// 	}
	// 	a.connected = false
	// }
	return
}

// Finalize finalizes the BLEAdaptor
func (b *BLEAdaptor) Finalize() (errs []error) {
	return b.Disconnect()
}

// ReadCharacteristic returns bytes from the BLE device for the 
// requested service and characteristic
func (b *BLEAdaptor) ReadCharacteristic(sUUID string, cUUID string) (data []byte, err error) {
	defer b.peripheral.Device().CancelConnection(b.peripheral)

	if !b.connected {
		log.Fatalf("Cannot read from BLE device until connected")
		return
	}
	
	c := make(chan []byte)
	f := func(gatt.Peripheral, error) {
		b.performRead(c, sUUID, cUUID)
	}

	b.device.Handle(
		gatt.PeripheralConnected(f),
	)

	b.peripheral.Device().Connect(b.peripheral)

	return <-c, nil
}

func (b *BLEAdaptor) performRead(c chan []byte, sUUID string, cUUID string) {
	s := b.getService(sUUID)
	characteristic := b.getCharacteristic(s, cUUID)

	val, err := b.peripheral.ReadCharacteristic(characteristic)
	if err != nil {
		fmt.Printf("Failed to read characteristic, err: %s\n", err)
		c <- []byte{}
	}

	c <- val
}

func (b *BLEAdaptor) getPeripheral() {

}

func (b *BLEAdaptor) getService(sUUID string) (s *gatt.Service) {
	// TODO: get the service that matches sUUID
	ss, err := b.peripheral.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	return ss[0] // TODO: return real one
}

func (b *BLEAdaptor) getCharacteristic(s *gatt.Service, cUUID string) (c *gatt.Characteristic) {
	// TODO: get the characteristic that matches cUUID
	cs, err := b.peripheral.DiscoverCharacteristics(nil, s)
	if err != nil {
		fmt.Printf("Failed to discover characteristics, err: %s\n", err)
		return
	}

	return cs[0] // TODO: return real one
}

func (b *BLEAdaptor) onStateChanged(d gatt.Device, s gatt.State) {
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

func (b *BLEAdaptor) onDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	fmt.Println("Discovered")
	id := strings.ToUpper(b.UUID())
	if strings.ToUpper(p.ID()) != id {
		return
	}

	b.connected = true
	b.peripheral = p

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("  Service Data      =", a.ServiceData)
	fmt.Println("")
}

func (b *BLEAdaptor) onConnected(p gatt.Peripheral, err error) {
	fmt.Println("Connected")
	defer p.Device().CancelConnection(p)
}

func (b *BLEAdaptor) onDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

