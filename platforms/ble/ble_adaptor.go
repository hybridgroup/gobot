package ble

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/paypal/gatt"
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
	name            string
	uuid            string
	device          gatt.Device
	peripheral      gatt.Peripheral
	services				map[string]*BLEService
	connected       bool
	ready	chan struct{}
	//connect   func(string) (io.ReadWriteCloser, error)
}

// NewBLEAdaptor returns a new BLEAdaptor given a name and uuid
func NewBLEAdaptor(name string, uuid string) *BLEAdaptor {
	return &BLEAdaptor{
		name:      name,
		uuid:      uuid,
		connected: false,
		ready: make(chan struct{}),
		services: make(map[string]*BLEService),
		// connect: func(port string) (io.ReadWriteCloser, error) {
		// 	return serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
		// },
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
		gatt.PeripheralDiscovered(b.onDiscovered),
		gatt.PeripheralConnected(b.onConnected),
		gatt.PeripheralDisconnected(b.onDisconnected),
	)

	device.Init(b.onStateChanged)
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
func (b *BLEAdaptor) ReadCharacteristic(sUUID string, cUUID string) (data chan []byte, err error) {
	//defer b.peripheral.Device().CancelConnection(b.peripheral)
	fmt.Println("ReadCharacteristic")
	if !b.connected {
		log.Fatalf("Cannot read from BLE device until connected")
		return
	}

	c := make(chan []byte)
	b.performRead(c, sUUID, cUUID)
	return c, nil
}

func (b *BLEAdaptor) performRead(c chan []byte, sUUID string, cUUID string) {
	fmt.Println("performRead")
	characteristic := b.services[sUUID].characteristics[cUUID]

	val, err := b.peripheral.ReadCharacteristic(characteristic)
	if err != nil {
		fmt.Printf("Failed to read characteristic, err: %s\n", err)
		c <- []byte{}
	}

	fmt.Printf("    value         %x | %q\n", val, val)
	c <- val
}

func (b *BLEAdaptor) getPeripheral() {

}

func (b *BLEAdaptor) getService(sUUID string) (service *gatt.Service) {
	fmt.Println("getService")
	ss, err := b.Peripheral().DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	fmt.Println("service")

	for _, s := range ss {
		msg := "Service: " + s.UUID().String()
		if len(s.Name()) > 0 {
			msg += " (" + s.Name() + ")"
		}
		fmt.Println(msg)

		id := strings.ToUpper(s.UUID().String())
		if strings.ToUpper(sUUID) != id {
			continue
		}

		msg = "Found Service: " + s.UUID().String()
		if len(s.Name()) > 0 {
			msg += " (" + s.Name() + ")"
		}
		fmt.Println(msg)
		return s
	}

	fmt.Println("getService: none found")
	return
}

func (b *BLEAdaptor) getCharacteristic(s *gatt.Service, cUUID string) (c *gatt.Characteristic) {
	fmt.Println("getCharacteristic")
	cs, err := b.Peripheral().DiscoverCharacteristics(nil, s)
	if err != nil {
		fmt.Printf("Failed to discover characteristics, err: %s\n", err)
		return
	}

	for _, char := range cs {
		id := strings.ToUpper(char.UUID().String())
		if strings.ToUpper(cUUID) != id {
			continue
		}

		msg := "  Found Characteristic  " + char.UUID().String()
		if len(char.Name()) > 0 {
			msg += " (" + char.Name() + ")"
		}
		msg += "\n    properties    " + char.Properties().String()
		fmt.Println(msg)
		return char
	}

	return nil
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
	id := strings.ToUpper(b.UUID())
	if strings.ToUpper(p.ID()) != id {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	// and connect to it
	p.Device().Connect(p)
}

func (b *BLEAdaptor) onConnected(p gatt.Peripheral, err error) {
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

func (b *BLEAdaptor) onDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

// Represents a BLE Peripheral's Service
type BLEService struct {
	uuid            	string
	service        		*gatt.Service
	characteristics 	map[string]*gatt.Characteristic
}

// NewBLEAdaptor returns a new BLEService given a uuid
func NewBLEService(sUuid string, service *gatt.Service) *BLEService {
	return &BLEService{
		uuid:      sUuid,
		service: 	 service,
		characteristics: make(map[string]*gatt.Characteristic),
	}
}
