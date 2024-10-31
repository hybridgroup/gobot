package tello

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot/v2"
)

const (
	// BounceEvent event
	BounceEvent = "bounce"

	// ConnectedEvent event
	ConnectedEvent = "connected"

	// FlightDataEvent event
	FlightDataEvent = "flightdata"

	// TakeoffEvent event
	TakeoffEvent = "takeoff"

	// LandingEvent event
	LandingEvent = "landing"

	// PalmLandingEvent event
	PalmLandingEvent = "palm-landing"

	// FlipEvent event
	FlipEvent = "flip"

	// TimeEvent event
	TimeEvent = "time"

	// LogEvent event
	LogEvent = "log"

	// WifiDataEvent event
	WifiDataEvent = "wifidata"

	// LightStrengthEvent event
	LightStrengthEvent = "lightstrength"

	// SetExposureEvent event
	SetExposureEvent = "setexposure"

	// VideoFrameEvent event
	VideoFrameEvent = "videoframe"

	// SetVideoEncoderRateEvent event
	SetVideoEncoderRateEvent = "setvideoencoder"
)

// the 16-bit messages and commands stored in bytes 6 & 5 of the packet
const (
	messageStart   = 0x00cc // 204
	wifiMessage    = 0x001a // 26
	videoRateQuery = 0x0028 // 40
	lightMessage   = 0x0035 // 53
	flightMessage  = 0x0056 // 86
	logMessage     = 0x1050 // 4176

	videoEncoderRateCommand = 0x0020 // 32
	videoStartCommand       = 0x0025 // 37
	exposureCommand         = 0x0034 // 52
	timeCommand             = 0x0046 // 70
	stickCommand            = 0x0050 // 80
	takeoffCommand          = 0x0054 // 84
	landCommand             = 0x0055 // 85
	flipCommand             = 0x005c // 92
	throwtakeoffCommand     = 0x005d // 93
	palmLandCommand         = 0x005e // 94
	bounceCommand           = 0x1053 // 4179
)

// FlipType is used for the various flips supported by the Tello.
type FlipType int

const (
	// FlipFront flips forward.
	FlipFront FlipType = 0

	// FlipLeft flips left.
	FlipLeft FlipType = 1

	// FlipBack flips backwards.
	FlipBack FlipType = 2

	// FlipRight flips to the right.
	FlipRight FlipType = 3

	// FlipForwardLeft flips forwards and to the left.
	FlipForwardLeft FlipType = 4

	// FlipBackLeft flips backwards and to the left.
	FlipBackLeft FlipType = 5

	// FlipBackRight flips backwards and to the right.
	FlipBackRight FlipType = 6

	// FlipForwardRight flips forewards and to the right.
	FlipForwardRight FlipType = 7
)

// VideoBitRate is used to set the bit rate for the streaming video returned by the Tello.
type VideoBitRate int

const (
	// VideoBitRateAuto sets the bitrate for streaming video to auto-adjust.
	VideoBitRateAuto VideoBitRate = 0

	// VideoBitRate1M sets the bitrate for streaming video to 1 Mb/s.
	VideoBitRate1M VideoBitRate = 1

	// VideoBitRate15M sets the bitrate for streaming video to 1.5 Mb/s
	VideoBitRate15M VideoBitRate = 2

	// VideoBitRate2M sets the bitrate for streaming video to 2 Mb/s.
	VideoBitRate2M VideoBitRate = 3

	// VideoBitRate3M sets the bitrate for streaming video to 3 Mb/s.
	VideoBitRate3M VideoBitRate = 4

	// VideoBitRate4M sets the bitrate for streaming video to 4 Mb/s.
	VideoBitRate4M VideoBitRate = 5
)

// FlightData packet returned by the Tello.
//
// The meaning of some fields is not documented. If you learned more, please, contribute.
// See https://github.com/hybridgroup/gobot/issues/798.
type FlightData struct {
	BatteryLow               bool
	BatteryLower             bool
	BatteryPercentage        int8 // How much battery left [in %].
	CameraState              int8
	DroneBatteryLeft         int16
	DroneFlyTimeLeft         int16
	DroneHover               bool // If the drone is in the air and not moving.
	EmOpen                   bool
	Flying                   bool  // If the drone is currently in the air.
	OnGround                 bool  // If the drone is currently on the ground.
	EastSpeed                int16 // Movement speed towards East [in cm/s]. Negative if moving west.
	ElectricalMachineryState int16
	FactoryMode              bool
	FlyMode                  int8
	FlyTime                  int16 // How long since take off [in s/10].
	FrontIn                  bool
	FrontLSC                 bool
	FrontOut                 bool
	GravityState             bool
	VerticalSpeed            int16 // Movement speed up [in cm/s].
	Height                   int16 // The height [in decimeters].
	ImuCalibrationState      int8  // The IMU calibration step (when doing IMU calibration).
	NorthSpeed               int16 // Movement speed towards North [in cm/s]. Negative if moving South.
	ThrowFlyTimer            int8

	// Warnings:
	DownVisualState bool // If the ground is visible by the down camera.
	BatteryState    bool // If there is an issue with battery.
	ImuState        bool // If drone needs IMU (Inertial Measurement Unit) calibration.
	OutageRecording bool // If there is an issue with video recording.
	PowerState      bool // If there is an issue with power supply.
	PressureState   bool // If there is an issue with air pressure.
	TemperatureHigh bool // If drone is overheating.
	WindState       bool // If the wind is too strong.
}

// WifiData packet returned by the Tello
type WifiData struct {
	Disturb  int8
	Strength int8
}

// Driver represents the DJI Tello drone
type Driver struct {
	name           string
	reqAddr        string
	cmdConn        io.WriteCloser // UDP connection to send/receive drone commands
	videoConn      *net.UDPConn   // UDP connection for drone video
	respPort       string
	videoPort      string
	cmdMutex       sync.Mutex
	seq            int16
	rx, ry, lx, ly float32
	throttle       int
	bouncing       bool
	gobot.Eventer
	doneCh            chan struct{}
	doneChReaderCount int32
}

// NewDriver creates a driver for the Tello drone. Pass in the UDP port to use for the responses
// from the drone.
func NewDriver(port string) *Driver {
	return NewDriverWithIP("192.168.10.1", port)
}

// NewDriverWithIP creates a driver for the Tello EDU drone. Pass in the ip address and UDP port to use for
// the responses from the drone.
func NewDriverWithIP(ip string, port string) *Driver {
	d := &Driver{
		name:      gobot.DefaultName("Tello"),
		reqAddr:   ip + ":8889",
		respPort:  port,
		videoPort: "11111",
		Eventer:   gobot.NewEventer(),
		doneCh:    make(chan struct{}, 1),
	}

	d.AddEvent(ConnectedEvent)
	d.AddEvent(FlightDataEvent)
	d.AddEvent(TakeoffEvent)
	d.AddEvent(LandingEvent)
	d.AddEvent(PalmLandingEvent)
	d.AddEvent(BounceEvent)
	d.AddEvent(FlipEvent)
	d.AddEvent(TimeEvent)
	d.AddEvent(LogEvent)
	d.AddEvent(WifiDataEvent)
	d.AddEvent(LightStrengthEvent)
	d.AddEvent(SetExposureEvent)
	d.AddEvent(VideoFrameEvent)
	d.AddEvent(SetVideoEncoderRateEvent)

	return d
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection { return nil }

// Start starts the driver.
func (d *Driver) Start() error {
	reqAddr, err := net.ResolveUDPAddr("udp", d.reqAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	respPort, err := net.ResolveUDPAddr("udp", ":"+d.respPort)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmdConn, err := net.DialUDP("udp", respPort, reqAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d.cmdConn = cmdConn

	// handle responses
	d.addDoneChReaderCount(1)
	go func() {
		defer d.addDoneChReaderCount(-1)

		err := d.On(d.Event(ConnectedEvent), func(interface{}) {
			if err := d.SendDateTime(); err != nil {
				panic(err)
			}
			if err := d.processVideo(); err != nil {
				panic(err)
			}
		})
		if err != nil {
			panic(err)
		}

	cmdLoop:
		for {
			select {
			case <-d.doneCh:
				break cmdLoop
			default:
				err := d.handleResponse(cmdConn)
				if err != nil {
					fmt.Println("response parse error:", err)
				}
			}
		}
	}()

	// starts notifications coming from drone to video port normally 11111
	if err := d.SendCommand(d.connectionString()); err != nil {
		return err
	}

	// send stick commands
	d.addDoneChReaderCount(1)
	go func() {
		defer d.addDoneChReaderCount(-1)

	stickCmdLoop:
		for {
			select {
			case <-d.doneCh:
				break stickCmdLoop
			default:
				if err := d.SendStickCommand(); err != nil {
					fmt.Println("stick command error:", err)
				}
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()

	return nil
}

// Halt stops the driver.
func (d *Driver) Halt() error {
	// send a landing command when we disconnect, and give it 500ms to be received before we shutdown
	if d.cmdConn != nil {
		if err := d.Land(); err != nil {
			return err
		}
	}
	time.Sleep(500 * time.Millisecond)

	if d.cmdConn != nil {
		d.cmdConn.Close()
	}

	if d.videoConn != nil {
		d.videoConn.Close()
	}
	readerCount := atomic.LoadInt32(&d.doneChReaderCount)
	for i := 0; i < int(readerCount); i++ {
		d.doneCh <- struct{}{}
	}

	return nil
}

// TakeOff tells drones to liftoff and start flying.
func (d *Driver) TakeOff() error {
	buf, _ := d.createPacket(takeoffCommand, 0x68, 0)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// Throw & Go support
func (d *Driver) ThrowTakeOff() error {
	buf, _ := d.createPacket(throwtakeoffCommand, 0x48, 0)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// Land tells drone to come in for landing.
func (d *Driver) Land() error {
	buf, _ := d.createPacket(landCommand, 0x68, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x00)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// StopLanding tells drone to stop landing.
func (d *Driver) StopLanding() error {
	buf, _ := d.createPacket(landCommand, 0x68, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x01)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// PalmLand tells drone to come in for a hand landing.
func (d *Driver) PalmLand() error {
	buf, _ := d.createPacket(palmLandCommand, 0x68, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x00)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// StartVideo tells Tello to send start info (SPS/PPS) for video stream.
func (d *Driver) StartVideo() error {
	buf, _ := d.createPacket(videoStartCommand, 0x60, 0)
	// seq = 0
	if err := binary.Write(buf, binary.LittleEndian, int16(0x00)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// SetExposure sets the drone camera exposure level. Valid levels are 0, 1, and 2.
func (d *Driver) SetExposure(level int) error {
	if level < 0 || level > 2 {
		return errors.New("Invalid exposure level")
	}

	buf, _ := d.createPacket(exposureCommand, 0x48, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(level)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// SetVideoEncoderRate sets the drone video encoder rate.
func (d *Driver) SetVideoEncoderRate(rate VideoBitRate) error {
	buf, _ := d.createPacket(videoEncoderRateCommand, 0x68, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(rate)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// SetFastMode sets the drone throttle to 1.
func (d *Driver) SetFastMode() error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.throttle = 1
	return nil
}

// SetSlowMode sets the drone throttle to 0.
func (d *Driver) SetSlowMode() error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.throttle = 0
	return nil
}

// Rate queries the current video bit rate.
func (d *Driver) Rate() error {
	buf, _ := d.createPacket(videoRateQuery, 0x48, 0)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// bound is a naive implementation that returns the smaller of x or y.
func bound(x, y float32) float32 { //nolint:unparam // keep y as parameter
	if x < -y {
		return -y
	}
	if x > y {
		return y
	}
	return x
}

// Vector returns the current motion vector.
// Values are from 0 to 1.
// x, y, z denote forward, side and vertical translation,
// and psi  yaw (rotation around the z-axis).
//
//nolint:nonamedreturns // sufficient here
func (d *Driver) Vector() (x, y, z, psi float32) {
	return d.ry, d.rx, d.ly, d.lx
}

// AddVector adds to the current motion vector.
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) AddVector(x, y, z, psi float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ry = bound(d.ry+x, 1)
	d.rx = bound(d.rx+y, 1)
	d.ly = bound(d.ly+z, 1)
	d.lx = bound(d.lx+psi, 1)

	return nil
}

// SetVector sets the current motion vector.
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) SetVector(x, y, z, psi float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ry = x
	d.rx = y
	d.ly = z
	d.lx = psi

	return nil
}

// SetX sets the x component of the current motion vector
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) SetX(x float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ry = x

	return nil
}

// SetY sets the y component of the current motion vector
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) SetY(y float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.rx = y

	return nil
}

// SetZ sets the z component of the current motion vector
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) SetZ(z float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ly = z

	return nil
}

// SetPsi sets the psi component (yaw) of the current motion vector
// Pass values from 0 to 1.
// See Vector() for the frame of reference.
func (d *Driver) SetPsi(psi float32) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.lx = psi

	return nil
}

// Up tells the drone to ascend. Pass in an int from 0-100.
func (d *Driver) Up(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ly = float32(val) / 100.0
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (d *Driver) Down(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ly = float32(val) / 100.0 * -1
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (d *Driver) Forward(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ry = float32(val) / 100.0
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (d *Driver) Backward(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.ry = float32(val) / 100.0 * -1
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (d *Driver) Right(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.rx = float32(val) / 100.0
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (d *Driver) Left(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.rx = float32(val) / 100.0 * -1
	return nil
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 0-100.
func (d *Driver) Clockwise(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.lx = float32(val) / 100.0
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 0-100.
func (d *Driver) CounterClockwise(val int) error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.lx = float32(val) / 100.0 * -1
	return nil
}

// Hover tells the drone to stop moving on the X, Y, and Z axes and stay in place
func (d *Driver) Hover() {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.rx = float32(0)
	d.ry = float32(0)
	d.lx = float32(0)
	d.ly = float32(0)
}

// CeaseRotation stops any rotational motion
func (d *Driver) CeaseRotation() {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	d.lx = float32(0)
}

// Bounce tells drone to start/stop performing the bouncing action
func (d *Driver) Bounce() error {
	buf, _ := d.createPacket(bounceCommand, 0x68, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}

	if d.bouncing {
		if err := binary.Write(buf, binary.LittleEndian, byte(0x31)); err != nil {
			return err
		}
	} else {
		if err := binary.Write(buf, binary.LittleEndian, byte(0x30)); err != nil {
			return err
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}
	_, err := d.cmdConn.Write(buf.Bytes())
	d.bouncing = !d.bouncing
	return err
}

// Flip tells drone to flip
func (d *Driver) Flip(direction FlipType) error {
	buf, _ := d.createPacket(flipCommand, 0x70, 1)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(direction)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// FrontFlip tells the drone to perform a front flip.
func (d *Driver) FrontFlip() error {
	return d.Flip(FlipFront)
}

// BackFlip tells the drone to perform a back flip.
func (d *Driver) BackFlip() error {
	return d.Flip(FlipBack)
}

// RightFlip tells the drone to perform a flip to the right.
func (d *Driver) RightFlip() error {
	return d.Flip(FlipRight)
}

// LeftFlip tells the drone to perform a flip to the left.
func (d *Driver) LeftFlip() error {
	return d.Flip(FlipLeft)
}

// ParseFlightData from drone
func (d *Driver) ParseFlightData(b []byte) (*FlightData, error) {
	buf := bytes.NewReader(b)
	fd := &FlightData{}
	var data byte

	if buf.Len() < 24 {
		err := errors.New("Invalid buffer length for flight data packet")
		fmt.Println(err)
		return fd, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &fd.Height); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.NorthSpeed); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.EastSpeed); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.VerticalSpeed); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.FlyTime); err != nil {
		return fd, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return fd, err
	}
	fd.ImuState = (data >> 0 & 0x1) == 1
	fd.PressureState = (data >> 1 & 0x1) == 1
	fd.DownVisualState = (data >> 2 & 0x1) == 1
	fd.PowerState = (data >> 3 & 0x1) == 1
	fd.BatteryState = (data >> 4 & 0x1) == 1
	fd.GravityState = (data >> 5 & 0x1) == 1
	fd.WindState = (data >> 7 & 0x1) == 1

	if err := binary.Read(buf, binary.LittleEndian, &fd.ImuCalibrationState); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.BatteryPercentage); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.DroneFlyTimeLeft); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.DroneBatteryLeft); err != nil {
		return fd, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return fd, err
	}
	fd.Flying = (data >> 0 & 0x1) == 1
	fd.OnGround = (data >> 1 & 0x1) == 1
	fd.EmOpen = (data >> 2 & 0x1) == 1
	fd.DroneHover = (data >> 3 & 0x1) == 1
	fd.OutageRecording = (data >> 4 & 0x1) == 1
	fd.BatteryLow = (data >> 5 & 0x1) == 1
	fd.BatteryLower = (data >> 6 & 0x1) == 1
	fd.FactoryMode = (data >> 7 & 0x1) == 1

	if err := binary.Read(buf, binary.LittleEndian, &fd.FlyMode); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.ThrowFlyTimer); err != nil {
		return fd, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &fd.CameraState); err != nil {
		return fd, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return fd, err
	}
	fd.ElectricalMachineryState = int16(data & 0xff)

	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return fd, err
	}
	fd.FrontIn = (data >> 0 & 0x1) == 1
	fd.FrontOut = (data >> 1 & 0x1) == 1
	fd.FrontLSC = (data >> 2 & 0x1) == 1

	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return fd, err
	}
	fd.TemperatureHigh = (data >> 0 & 0x1) == 1

	return fd, nil
}

// SendStickCommand sends the joystick command packet to the drone.
func (d *Driver) SendStickCommand() error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	buf, _ := d.createPacket(stickCommand, 0x60, 11)
	// seq = 0
	if err := binary.Write(buf, binary.LittleEndian, int16(0x00)); err != nil {
		return err
	}

	// RightX center=1024 left =364 right =-364
	axis1 := int16(660.0*d.rx + 1024.0)

	// RightY down =364 up =-364
	axis2 := int16(660.0*d.ry + 1024.0)

	// LeftY down =364 up =-364
	axis3 := int16(660.0*d.ly + 1024.0)

	// LeftX left =364 right =-364
	axis4 := int16(660.0*d.lx + 1024.0)

	// speed control
	axis5 := int16(d.throttle) //nolint:gosec // TODO: fix later

	packedAxis := int64(axis1)&0x7FF | int64(axis2&0x7FF)<<11 | 0x7FF&int64(axis3)<<22 | 0x7FF&int64(axis4)<<33 |
		int64(axis5)<<44
	if err := binary.Write(buf, binary.LittleEndian, byte(0xFF&packedAxis)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(packedAxis>>8&0xFF)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(packedAxis>>16&0xFF)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(packedAxis>>24&0xFF)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(packedAxis>>32&0xFF)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(packedAxis>>40&0xFF)); err != nil {
		return err
	}

	now := time.Now()
	if err := binary.Write(buf, binary.LittleEndian, byte(now.Hour())); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(now.Minute())); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(now.Second())); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(now.UnixNano()/int64(time.Millisecond)&0xff)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(now.UnixNano()/int64(time.Millisecond)>>8)); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())

	return err
}

// SendDateTime sends the current date/time to the drone.
func (d *Driver) SendDateTime() error {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	buf, _ := d.createPacket(timeCommand, 0x50, 11)
	d.seq++
	if err := binary.Write(buf, binary.LittleEndian, d.seq); err != nil {
		return err
	}

	now := time.Now()
	if err := binary.Write(buf, binary.LittleEndian, byte(0x00)); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(buf, binary.LittleEndian, int16(now.Hour())); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(buf, binary.LittleEndian, int16(now.Minute())); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(buf, binary.LittleEndian, int16(now.Second())); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(buf, binary.LittleEndian, int16(now.UnixNano()/int64(time.Millisecond)&0xff)); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(buf, binary.LittleEndian, int16(now.UnixNano()/int64(time.Millisecond)>>8)); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes())); err != nil {
		return err
	}

	_, err := d.cmdConn.Write(buf.Bytes())
	return err
}

// SendCommand is used to send a text command such as the initial connection request to the drone.
func (d *Driver) SendCommand(cmd string) error {
	_, err := d.cmdConn.Write([]byte(cmd))
	return err
}

func (d *Driver) handleResponse(r io.Reader) error {
	var buf [2048]byte
	var msgType uint16
	n, err := r.Read(buf[0:])
	if err != nil {
		return err
	}

	// parse binary packet
	if buf[0] == messageStart {
		msgType = (uint16(buf[6]) << 8) | uint16(buf[5])
		switch msgType {
		case wifiMessage:
			wd := &WifiData{
				Strength: int8(buf[9:10][0]),
				Disturb:  int8(buf[10:11][0]),
			}
			d.Publish(d.Event(WifiDataEvent), wd)
		case lightMessage:
			d.Publish(d.Event(LightStrengthEvent), int8(buf[9:10][0]))
		case logMessage:
			d.Publish(d.Event(LogEvent), buf[9:])
		case timeCommand:
			d.Publish(d.Event(TimeEvent), buf[7:8])
		case bounceCommand:
			d.Publish(d.Event(BounceEvent), buf[7:8])
		case takeoffCommand:
			d.Publish(d.Event(TakeoffEvent), buf[7:8])
		case landCommand:
			d.Publish(d.Event(LandingEvent), buf[7:8])
		case palmLandCommand:
			d.Publish(d.Event(PalmLandingEvent), buf[7:8])
		case flipCommand:
			d.Publish(d.Event(FlipEvent), buf[7:8])
		case flightMessage:
			fd, _ := d.ParseFlightData(buf[9:])
			d.Publish(d.Event(FlightDataEvent), fd)
		case exposureCommand:
			d.Publish(d.Event(SetExposureEvent), buf[7:8])
		case videoEncoderRateCommand:
			d.Publish(d.Event(SetVideoEncoderRateEvent), buf[7:8])
		default:
			fmt.Printf("Unknown message: %+v\n", buf[0:n])
		}
		return nil
	}

	// parse text packet
	if buf[0] == 0x63 && buf[1] == 0x6f && buf[2] == 0x6e {
		d.Publish(d.Event(ConnectedEvent), nil)
	}

	return nil
}

func (d *Driver) processVideo() error {
	videoPort, err := net.ResolveUDPAddr("udp", ":11111")
	if err != nil {
		return err
	}
	d.videoConn, err = net.ListenUDP("udp", videoPort)
	if err != nil {
		return err
	}

	d.addDoneChReaderCount(1)
	go func() {
		defer d.addDoneChReaderCount(-1)

	videoConnLoop:
		for {
			select {
			case <-d.doneCh:
				break videoConnLoop
			default:
				buf := make([]byte, 2048)
				n, _, err := d.videoConn.ReadFromUDP(buf)
				if err != nil {
					fmt.Println("Error: ", err)
					continue
				}

				d.Publish(d.Event(VideoFrameEvent), buf[2:n])
			}
		}
	}()

	return nil
}

func (d *Driver) createPacket(cmd int16, pktType byte, pktLen int16) (*bytes.Buffer, error) {
	l := pktLen + 11
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, binary.LittleEndian, byte(messageStart)); err != nil {
		return buf, err
	}
	if err := binary.Write(buf, binary.LittleEndian, l<<3); err != nil {
		return buf, err
	}
	if err := binary.Write(buf, binary.LittleEndian, CalculateCRC8(buf.Bytes()[0:3])); err != nil {
		return buf, err
	}
	if err := binary.Write(buf, binary.LittleEndian, pktType); err != nil {
		return buf, err
	}
	if err := binary.Write(buf, binary.LittleEndian, cmd); err != nil {
		return buf, err
	}

	return buf, nil
}

func (d *Driver) connectionString() string {
	x, _ := strconv.Atoi(d.videoPort)
	b := [2]byte{}
	binary.LittleEndian.PutUint16(b[:], uint16(x)) //nolint:gosec // TODO: fix later
	res := fmt.Sprintf("conn_req:%s", b)
	return res
}

func (d *Driver) addDoneChReaderCount(delta int32) {
	atomic.AddInt32(&d.doneChReaderCount, delta)
}

func (f *FlightData) AirSpeed() float64 {
	return math.Sqrt(
		math.Pow(float64(f.NorthSpeed), 2) +
			math.Pow(float64(f.EastSpeed), 2) +
			math.Pow(float64(f.VerticalSpeed), 2))
}

func (f *FlightData) GroundSpeed() float64 {
	return math.Sqrt(
		math.Pow(float64(f.NorthSpeed), 2) +
			math.Pow(float64(f.EastSpeed), 2))
}
