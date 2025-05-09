package parrot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// droneCommandService      = "9a66fa000800919111e4012d1540cb8e"
	// droneNotificationService = "9a66fb000800919111e4012d1540cb8e"

	// send characteristics
	pcmdChara     = "9a66fa0a0800919111e4012d1540cb8e"
	commandChara  = "9a66fa0b0800919111e4012d1540cb8e"
	priorityChara = "9a66fa0c0800919111e4012d1540cb8e"

	// receive characteristics
	flightStatusChara = "9a66fb0e0800919111e4012d1540cb8e"
	batteryChara      = "9a66fb0f0800919111e4012d1540cb8e"

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

	BatteryEvent        = "battery"
	FlightStatusEvent   = "flightstatus"
	TakeoffEvent        = "takeoff"
	HoveringEvent       = "hovering"
	FlyingEvent         = "flying"
	LandingEvent        = "landing"
	LandedEvent         = "landed"
	EmergencyEvent      = "emergency"
	RollingEvent        = "rolling"
	FlatTrimChangeEvent = "flattrimchange"

	// modes for LightControl
	LightFixed      = 0
	LightBlinked    = 1
	LightOscillated = 3

	// modes for ClawControl
	ClawOpen   = 0
	ClawClosed = 1
)

// MinidroneDriver is the Gobot interface to the Parrot Minidrone
type MinidroneDriver struct {
	*ble.Driver
	stepsfa0a uint16
	stepsfa0b uint16
	pcmdMutex sync.Mutex
	flying    bool
	Pcmd      Pcmd
	gobot.Eventer
}

// Pcmd is the Parrot Command structure for flight control
type Pcmd struct {
	Flag  int
	Roll  int
	Pitch int
	Yaw   int
	Gaz   int
	Psi   float32
}

// NewMinidroneDriver creates a Parrot Minidrone Driver
func NewMinidroneDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *MinidroneDriver {
	d := &MinidroneDriver{
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
	d.Driver = ble.NewDriver(a, "Minidrone", d.initialize, d.shutdown, opts...)

	d.AddEvent(BatteryEvent)
	d.AddEvent(FlightStatusEvent)
	d.AddEvent(TakeoffEvent)
	d.AddEvent(FlyingEvent)
	d.AddEvent(HoveringEvent)
	d.AddEvent(LandingEvent)
	d.AddEvent(LandedEvent)
	d.AddEvent(EmergencyEvent)
	d.AddEvent(RollingEvent)

	return d
}

// GenerateAllStates sets up all the default states aka settings on the drone
func (d *MinidroneDriver) GenerateAllStates() error {
	d.stepsfa0b++
	buf := []byte{
		0x04, byte(d.stepsfa0b), 0x00, 0x04, 0x01, 0x00, 0x32, 0x30, 0x31, 0x34, 0x2D, 0x31, 0x30, 0x2D, 0x32, 0x38, 0x00,
	}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// TakeOff tells the Minidrone to takeoff
func (d *MinidroneDriver) TakeOff() error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x00, 0x01, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// Land tells the Minidrone to land
func (d *MinidroneDriver) Land() error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x00, 0x03, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// FlatTrim calibrates the Minidrone to use its current position as being level
func (d *MinidroneDriver) FlatTrim() error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x00, 0x00, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// Emergency sets the Minidrone into emergency mode
func (d *MinidroneDriver) Emergency() error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x00, 0x04, 0x00}
	return d.Adaptor().WriteCharacteristic(priorityChara, buf)
}

// TakePicture tells the Minidrone to take a picture
func (d *MinidroneDriver) TakePicture() error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x06, 0x01, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// StartPcmd starts the continuous Pcmd communication with the Minidrone
func (d *MinidroneDriver) StartPcmd() {
	go func() {
		// wait a little bit so that there is enough time to get some ACKs
		time.Sleep(500 * time.Millisecond)
		for {
			err := d.Adaptor().WriteCharacteristic(pcmdChara, d.generatePcmd().Bytes())
			if err != nil {
				fmt.Println("pcmd write error:", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

// Up tells the drone to ascend. Pass in an int from 0-100.
func (d *MinidroneDriver) Up(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Gaz = validatePitch(val)
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (d *MinidroneDriver) Down(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Gaz = validatePitch(val) * -1
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (d *MinidroneDriver) Forward(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Pitch = validatePitch(val)
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (d *MinidroneDriver) Backward(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Pitch = validatePitch(val) * -1
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (d *MinidroneDriver) Right(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Roll = validatePitch(val)
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (d *MinidroneDriver) Left(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Roll = validatePitch(val) * -1
	return nil
}

// Clockwise tells drone to rotate in a clockwise directiod. Pass in an int from 0-100.
func (d *MinidroneDriver) Clockwise(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Yaw = validatePitch(val)
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise directiod.
// Pass in an int from 0-100.
func (d *MinidroneDriver) CounterClockwise(val int) error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd.Flag = 1
	d.Pcmd.Yaw = validatePitch(val) * -1
	return nil
}

// Stop tells the drone to stop moving in any direction and simply hover in place
func (d *MinidroneDriver) Stop() error {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()

	d.Pcmd = Pcmd{
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
func (d *MinidroneDriver) StartRecording() error {
	return nil
}

// StopRecording is not supported by the Parrot Minidrone
func (d *MinidroneDriver) StopRecording() error {
	return nil
}

// HullProtection is not supported by the Parrot Minidrone
func (d *MinidroneDriver) HullProtection(protect bool) error {
	return nil
}

// Outdoor mode is not supported by the Parrot Minidrone
func (d *MinidroneDriver) Outdoor(outdoor bool) error {
	return nil
}

// FrontFlip tells the drone to perform a front flip
func (d *MinidroneDriver) FrontFlip() error {
	return d.Adaptor().WriteCharacteristic(commandChara, d.generateAnimation(0).Bytes())
}

// BackFlip tells the drone to perform a backflip
func (d *MinidroneDriver) BackFlip() error {
	return d.Adaptor().WriteCharacteristic(commandChara, d.generateAnimation(1).Bytes())
}

// RightFlip tells the drone to perform a flip to the right
func (d *MinidroneDriver) RightFlip() error {
	return d.Adaptor().WriteCharacteristic(commandChara, d.generateAnimation(2).Bytes())
}

// LeftFlip tells the drone to perform a flip to the left
func (d *MinidroneDriver) LeftFlip() error {
	return d.Adaptor().WriteCharacteristic(commandChara, d.generateAnimation(3).Bytes())
}

// LightControl controls lights on those Minidrone models which
// have the correct hardware, such as the Maclane, Blaze, & Swat.
// Params:
//
//	id - always 0
//	mode - either LightFixed, LightBlinked, or LightOscillated
//	intensity - Light intensity from 0 (OFF) to 100 (Max intensity).
//				Only used in LightFixed mode.
func (d *MinidroneDriver) LightControl(id uint8, mode uint8, intensity uint8) error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x10, 0x00, id, mode, intensity, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// ClawControl controls the claw on the Parrot Mambo
// Params:
//
//	id - always 0
//	mode - either ClawOpen or ClawClosed
func (d *MinidroneDriver) ClawControl(id uint8, mode uint8) error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x10, 0x01, id, mode, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// GunControl fires the gun on the Parrot Mambo
// Params:
//
//	id - always 0
func (d *MinidroneDriver) GunControl(id uint8) error {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x10, 0x02, id, 0x00}
	return d.Adaptor().WriteCharacteristic(commandChara, buf)
}

// initialize tells driver to get ready to do work
func (d *MinidroneDriver) initialize() error {
	d.Adaptor().WithoutResponses(true)

	if err := d.GenerateAllStates(); err != nil {
		return err
	}

	// subscribe to battery notifications
	if err := d.Adaptor().Subscribe(batteryChara, func(data []byte) {
		d.Publish(d.Event(BatteryEvent), data[len(data)-1])
	}); err != nil {
		return err
	}

	// subscribe to flying status notifications
	if err := d.Adaptor().Subscribe(flightStatusChara, func(data []byte) {
		d.processFlightStatus(data)
	}); err != nil {
		return err
	}

	if err := d.FlatTrim(); err != nil {
		return err
	}

	d.StartPcmd()

	return d.FlatTrim()
}

// shutdown stops minidrone driver (void)
func (d *MinidroneDriver) shutdown() error {
	err := d.Land()
	time.Sleep(500 * time.Millisecond)
	return err
}

func (d *MinidroneDriver) generateAnimation(direction int8) *bytes.Buffer {
	d.stepsfa0b++
	buf := []byte{0x02, byte(d.stepsfa0b) & 0xff, 0x02, 0x04, 0x00, 0x00, byte(direction), 0x00, 0x00, 0x00}
	return bytes.NewBuffer(buf)
}

func (d *MinidroneDriver) generatePcmd() *bytes.Buffer {
	d.pcmdMutex.Lock()
	defer d.pcmdMutex.Unlock()
	d.stepsfa0a++
	pcmd := d.Pcmd

	cmd := &bytes.Buffer{}
	if err := binary.Write(cmd, binary.LittleEndian, int8(2)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(d.stepsfa0a)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int8(2)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int8(0)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int8(2)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int8(0)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(pcmd.Flag)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(pcmd.Roll)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(pcmd.Pitch)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(pcmd.Yaw)); err != nil {
		panic(err)
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(cmd, binary.LittleEndian, int8(pcmd.Gaz)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, pcmd.Psi); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int16(0)); err != nil {
		panic(err)
	}
	if err := binary.Write(cmd, binary.LittleEndian, int16(0)); err != nil {
		panic(err)
	}

	return cmd
}

func (d *MinidroneDriver) processFlightStatus(data []byte) {
	if len(data) < 5 {
		// ignore, just a sync
		return
	}

	d.Publish(FlightStatusEvent, data[4])

	switch data[4] {
	case flatTrimChanged:
		d.Publish(FlatTrimChangeEvent, true)

	case flyingStateChanged:
		switch data[6] {
		case flyingStateLanded:
			if d.flying {
				d.flying = false
				d.Publish(LandedEvent, true)
			}
		case flyingStateTakeoff:
			d.Publish(TakeoffEvent, true)
		case flyingStateHovering:
			if !d.flying {
				d.flying = true
				d.Publish(HoveringEvent, true)
			}
		case flyingStateFlying:
			if !d.flying {
				d.flying = true
				d.Publish(FlyingEvent, true)
			}
		case flyingStateLanding:
			d.Publish(LandingEvent, true)
		case flyingStateEmergency:
			d.Publish(EmergencyEvent, true)
		case flyingStateRolling:
			d.Publish(RollingEvent, true)
		}
	}
}

// ValidatePitch helps validate pitch values such as those created by
// a joystick to values between 0-100 that are required as
// params to Parrot Minidrone PCMDs
func ValidatePitch(data float64, offset float64) int {
	value := math.Abs(data) / offset
	if value >= 0.1 {
		if value <= 1.0 {
			return int((float64(int(value*100)) / 100) * 100)
		}
		return 100
	}
	return 0
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}
