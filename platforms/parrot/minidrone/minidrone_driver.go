package minidrone

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// Driver is the Gobot interface to the Parrot Minidrone
type Driver struct {
	name       string
	connection gobot.Connection
	stepsfa0a  uint16
	stepsfa0b  uint16
	pcmdMutex  sync.Mutex
	flying     bool
	Pcmd       Pcmd
	gobot.Eventer
}

const (
	// BLE services
	droneCommandService      = "9a66fa000800919111e4012d1540cb8e"
	droneNotificationService = "9a66fb000800919111e4012d1540cb8e"

	// send characteristics
	pcmdCharacteristic     = "9a66fa0a0800919111e4012d1540cb8e"
	commandCharacteristic  = "9a66fa0b0800919111e4012d1540cb8e"
	priorityCharacteristic = "9a66fa0c0800919111e4012d1540cb8e"

	// receive characteristics
	flightStatusCharacteristic = "9a66fb0e0800919111e4012d1540cb8e"
	batteryCharacteristic      = "9a66fb0f0800919111e4012d1540cb8e"

	// piloting states
	flatTrimChanged    = 0
	flyingStateChanged = 1

	// flying states
	flyingStateLanded    = 0
	flyingStateTakeoff   = 1
	flyingStateHovering  = 2
	flyingStateFlying    = 3
	flyingStateLanding   = 4
	flyingStateEmergency = 5
	flyingStateRolling   = 6

	// Battery event
	Battery = "battery"

	// FlightStatus event
	FlightStatus = "flightstatus"

	// Takeoff event
	Takeoff = "takeoff"

	// Hovering event
	Hovering = "hovering"

	// Flying event
	Flying = "flying"

	// Landing event
	Landing = "landing"

	// Landed event
	Landed = "landed"

	// Emergency event
	Emergency = "emergency"

	// Rolling event
	Rolling = "rolling"

	// FlatTrimChange event
	FlatTrimChange = "flattrimchange"

	// LightFixed mode for LightControl
	LightFixed = 0

	// LightBlinked mode for LightControl
	LightBlinked = 1

	// LightOscillated mode for LightControl
	LightOscillated = 3

	// ClawOpen mode for ClawControl
	ClawOpen = 0

	// ClawClosed mode for ClawControl
	ClawClosed = 1
)

// Pcmd is the Parrot Command structure for flight control
type Pcmd struct {
	Flag  int
	Roll  int
	Pitch int
	Yaw   int
	Gaz   int
	Psi   float32
}

// NewDriver creates a Parrot Minidrone Driver
func NewDriver(a ble.BLEConnector) *Driver {
	n := &Driver{
		name:       gobot.DefaultName("Minidrone"),
		connection: a,
		Pcmd: Pcmd{
			Flag:  0,
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Gaz:   0,
			Psi:   0,
		},
		Eventer: gobot.NewEventer(),
	}

	n.AddEvent(Battery)
	n.AddEvent(FlightStatus)

	n.AddEvent(Takeoff)
	n.AddEvent(Flying)
	n.AddEvent(Hovering)
	n.AddEvent(Landing)
	n.AddEvent(Landed)
	n.AddEvent(Emergency)
	n.AddEvent(Rolling)

	return n
}

// Connection returns the BLE connection
func (b *Driver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *Driver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *Driver) Start() (err error) {
	b.adaptor().WithoutReponses(true)
	b.Init()
	b.FlatTrim()
	b.StartPcmd()
	b.FlatTrim()

	return
}

// Halt stops minidrone driver (void)
func (b *Driver) Halt() (err error) {
	b.Land()

	time.Sleep(500 * time.Millisecond)
	return
}

// Init initializes the BLE insterfaces used by the Minidrone
func (b *Driver) Init() (err error) {
	b.GenerateAllStates()

	// subscribe to battery notifications
	b.adaptor().Subscribe(batteryCharacteristic, func(data []byte, e error) {
		b.Publish(b.Event(Battery), data[len(data)-1])
	})

	// subscribe to flying status notifications
	b.adaptor().Subscribe(flightStatusCharacteristic, func(data []byte, e error) {
		b.processFlightStatus(data)
	})

	return
}

// GenerateAllStates sets up all the default states aka settings on the drone
func (b *Driver) GenerateAllStates() (err error) {
	b.stepsfa0b++
	buf := []byte{0x04, byte(b.stepsfa0b), 0x00, 0x04, 0x01, 0x00, 0x32, 0x30, 0x31, 0x34, 0x2D, 0x31, 0x30, 0x2D, 0x32, 0x38, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)

	return
}

// TakeOff tells the Minidrone to takeoff
func (b *Driver) TakeOff() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x01, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)

	return
}

// Land tells the Minidrone to land
func (b *Driver) Land() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x03, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)

	return err
}

// FlatTrim calibrates the Minidrone to use its current position as being level
func (b *Driver) FlatTrim() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x00, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)

	return err
}

// Emergency sets the Minidrone into emergency mode
func (b *Driver) Emergency() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x04, 0x00}
	err = b.adaptor().WriteCharacteristic(priorityCharacteristic, buf)

	return err
}

// TakePicture tells the Minidrone to take a picture
func (b *Driver) TakePicture() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x06, 0x01, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)

	return err
}

// StartPcmd starts the continuous Pcmd communication with the Minidrone
func (b *Driver) StartPcmd() {
	go func() {
		// wait a little bit so that there is enough time to get some ACKs
		time.Sleep(500 * time.Millisecond)
		for {
			err := b.adaptor().WriteCharacteristic(pcmdCharacteristic, b.generatePcmd().Bytes())
			if err != nil {
				fmt.Println("pcmd write error:", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

// Up tells the drone to ascend. Pass in an int from 0-100.
func (b *Driver) Up(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val)
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (b *Driver) Down(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val) * -1
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (b *Driver) Forward(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val)
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (b *Driver) Backward(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val) * -1
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (b *Driver) Right(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val)
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (b *Driver) Left(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val) * -1
	return nil
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 0-100.
func (b *Driver) Clockwise(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val)
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 0-100.
func (b *Driver) CounterClockwise(val int) error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val) * -1
	return nil
}

// Stop tells the drone to stop moving in any direction and simply hover in place
func (b *Driver) Stop() error {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()

	b.Pcmd = Pcmd{
		Flag:  0,
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
		Gaz:   0,
		Psi:   0,
	}

	return nil
}

// StartRecording is not supported by the Parrot Minidrone
func (b *Driver) StartRecording() error {
	return nil
}

// StopRecording is not supported by the Parrot Minidrone
func (b *Driver) StopRecording() error {
	return nil
}

// HullProtection is not supported by the Parrot Minidrone
func (b *Driver) HullProtection(protect bool) error {
	return nil
}

// Outdoor mode is not supported by the Parrot Minidrone
func (b *Driver) Outdoor(outdoor bool) error {
	return nil
}

// FrontFlip tells the drone to perform a front flip
func (b *Driver) FrontFlip() (err error) {
	return b.adaptor().WriteCharacteristic(commandCharacteristic, b.generateAnimation(0).Bytes())
}

// BackFlip tells the drone to perform a backflip
func (b *Driver) BackFlip() (err error) {
	return b.adaptor().WriteCharacteristic(commandCharacteristic, b.generateAnimation(1).Bytes())
}

// RightFlip tells the drone to perform a flip to the right
func (b *Driver) RightFlip() (err error) {
	return b.adaptor().WriteCharacteristic(commandCharacteristic, b.generateAnimation(2).Bytes())
}

// LeftFlip tells the drone to perform a flip to the left
func (b *Driver) LeftFlip() (err error) {
	return b.adaptor().WriteCharacteristic(commandCharacteristic, b.generateAnimation(3).Bytes())
}

// LightControl controls lights on those Minidrone models which
// have the correct hardware, such as the Maclane, Blaze, & Swat.
// Params:
//		id - always 0
//		mode - either LightFixed, LightBlinked, or LightOscillated
//		intensity - Light intensity from 0 (OFF) to 100 (Max intensity).
// 					Only used in LightFixed mode.
//
func (b *Driver) LightControl(id uint8, mode uint8, intensity uint8) (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x10, 0x00, id, mode, intensity, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)
	return
}

// ClawControl controls the claw on the Parrot Mambo
// Params:
//		id - always 0
//		mode - either ClawOpen or ClawClosed
//
func (b *Driver) ClawControl(id uint8, mode uint8) (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x10, 0x01, id, mode, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)
	return
}

// GunControl fires the gun on the Parrot Mambo
// Params:
//		id - always 0
//
func (b *Driver) GunControl(id uint8) (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x10, 0x02, id, 0x00}
	err = b.adaptor().WriteCharacteristic(commandCharacteristic, buf)
	return
}

func (b *Driver) generateAnimation(direction int8) *bytes.Buffer {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x04, 0x00, 0x00, byte(direction), 0x00, 0x00, 0x00}
	return bytes.NewBuffer(buf)
}

func (b *Driver) generatePcmd() *bytes.Buffer {
	b.pcmdMutex.Lock()
	defer b.pcmdMutex.Unlock()
	b.stepsfa0a++
	pcmd := b.Pcmd

	cmd := &bytes.Buffer{}
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(b.stepsfa0a))
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(0))
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(0))
	binary.Write(cmd, binary.LittleEndian, int8(pcmd.Flag))
	binary.Write(cmd, binary.LittleEndian, int8(pcmd.Roll))
	binary.Write(cmd, binary.LittleEndian, int8(pcmd.Pitch))
	binary.Write(cmd, binary.LittleEndian, int8(pcmd.Yaw))
	binary.Write(cmd, binary.LittleEndian, int8(pcmd.Gaz))
	binary.Write(cmd, binary.LittleEndian, float32(pcmd.Psi))
	binary.Write(cmd, binary.LittleEndian, int16(0))
	binary.Write(cmd, binary.LittleEndian, int16(0))

	return cmd
}

func (b *Driver) processFlightStatus(data []byte) {
	if len(data) < 5 {
		// ignore, just a sync
		return
	}

	b.Publish(FlightStatus, data[4])

	switch data[4] {
	case flatTrimChanged:
		b.Publish(FlatTrimChange, true)

	case flyingStateChanged:
		switch data[6] {
		case flyingStateLanded:
			if b.flying {
				b.flying = false
				b.Publish(Landed, true)
			}
		case flyingStateTakeoff:
			b.Publish(Takeoff, true)
		case flyingStateHovering:
			if !b.flying {
				b.flying = true
				b.Publish(Hovering, true)
			}
		case flyingStateFlying:
			if !b.flying {
				b.flying = true
				b.Publish(Flying, true)
			}
		case flyingStateLanding:
			b.Publish(Landing, true)
		case flyingStateEmergency:
			b.Publish(Emergency, true)
		case flyingStateRolling:
			b.Publish(Rolling, true)
		}
	}
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}
