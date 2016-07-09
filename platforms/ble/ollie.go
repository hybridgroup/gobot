package ble

import (
	"bytes"
	"fmt"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*SpheroOllieDriver)(nil)

type SpheroOllieDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

const (
	// service IDs
	SpheroBLEService = "22bb746f2bb075542d6f726568705327"

	// characteristic IDs
	WakeCharacteristic = "22bb746f2bbf75542d6f726568705327"
	TXPowerCharacteristic = "22bb746f2bb275542d6f726568705327"
	AntiDosCharacteristic = "22bb746f2bbd75542d6f726568705327"
	RobotControlService = "22bb746f2ba075542d6f726568705327"
	CommandsCharacteristic = "22bb746f2ba175542d6f726568705327"
	ResponseCharacteristic = "22bb746f2ba675542d6f726568705327"
)

// NewSpheroOllieDriver creates a SpheroOllieDriver by name
func NewSpheroOllieDriver(a *BLEClientAdaptor, name string) *SpheroOllieDriver {
	n := &SpheroOllieDriver{
		name:       name,
		connection: a,
		Eventer: gobot.NewEventer(),
	}

	return n
}
func (b *SpheroOllieDriver) Connection() gobot.Connection { return b.connection }
func (b *SpheroOllieDriver) Name() string                 { return b.name }

// adaptor returns BLE adaptor
func (b *SpheroOllieDriver) adaptor() *BLEClientAdaptor {
	return b.Connection().(*BLEClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *SpheroOllieDriver) Start() (errs []error) {
	b.Init()

	return
}

// Halt stops Ollie driver (void)
func (b *SpheroOllieDriver) Halt() (errs []error) {
	return
}

func (b *SpheroOllieDriver) Init() (err error) {
	b.AntiDOSOff()
	b.SetTXPower(7)
	b.Wake()

	// subscribe to Sphero response notifications
	b.adaptor().Subscribe(RobotControlService, ResponseCharacteristic, b.HandleResponses)

	return
}

// Turns off Anti-DOS code so we can control Ollie
func (b *SpheroOllieDriver) AntiDOSOff() (err error) {
	str := "011i3"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, AntiDosCharacteristic, buf.Bytes())
	if err != nil {
		fmt.Println("AntiDOSOff error:", err)
		return err
	}

	return
}

// Wakes Ollie up so we can play
func (b *SpheroOllieDriver) Wake() (err error) {
	buf := []byte{0x01}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, WakeCharacteristic, buf)
	if err != nil {
		fmt.Println("Wake error:", err)
		return err
	}

	return
}

// Sets transmit level
func (b *SpheroOllieDriver) SetTXPower(level int) (err error) {
	buf := []byte{byte(level)}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, TXPowerCharacteristic, buf)
	if err != nil {
		fmt.Println("SetTXLevel error:", err)
		return err
	}

	return
}

// Handle responses returned from Ollie
func (b *SpheroOllieDriver) HandleResponses(data []byte, e error) {
	fmt.Println("response data:", data)

	return
}

// SetRGB sets the Ollie to the given r, g, and b values
func (s *SpheroOllieDriver) SetRGB(r uint8, g uint8, b uint8) {
	fmt.Println("setrgb")
}

// Tells the Ollie to roll
func (s *SpheroOllieDriver) Roll(speed uint8, heading uint16) {
	fmt.Println("roll", speed, heading)
}

// Tells the Ollie to stop
func (s *SpheroOllieDriver) Stop() {
	s.Roll(0, 0)
}
