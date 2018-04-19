package tello

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

const (
	// ConnectedEvent event
	ConnectedEvent = "connected"

	// FlightDataEvent event
	FlightDataEvent = "flightdata"

	// TakeoffEvent event
	TakeoffEvent = "takeoff"

	// LandingEvent event
	LandingEvent = "landing"

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
)

const (
	messageStart  = 0xcc
	wifiMessage   = 26
	lightMessage  = 53
	timeMessage   = 70
	flightMessage = 86

	logMessage = 0x50

	videoEncoderRateCommand = 0x20
	videoStartCommand       = 0x25
	exposureCommand         = 0x34
	stickCommand            = 0x50
	takeoffCommand          = 0x54
	landCommand             = 0x55
	flipCommand             = 0x5c

	flipFront = 0
	flipLeft  = 1
	flipBack  = 2
	flipRight = 3
)

// FlightData packet returned by the Tello
type FlightData struct {
	batteryLow               int16
	batteryLower             int16
	batteryPercentage        int8
	batteryState             int16
	cameraState              int8
	downVisualState          int16
	droneBatteryLeft         int16
	droneFlyTimeLeft         int16
	droneHover               int16
	eMOpen                   int16
	eMSky                    int16
	eMGround                 int16
	eastSpeed                int16
	electricalMachineryState int16
	factoryMode              int16
	flyMode                  int8
	flySpeed                 int16
	flyTime                  int16
	frontIn                  int16
	frontLSC                 int16
	frontOut                 int16
	gravityState             int16
	groundSpeed              int16
	height                   int16
	imuCalibrationState      int8
	imuState                 int16
	lightStrength            int16
	northSpeed               int16
	outageRecording          int16
	powerState               int16
	pressureState            int16
	smartVideoExitMode       int16
	temperatureHeight        int16
	throwFlyTimer            int8
	wifiDisturb              int16
	wifiStrength             int16
	windState                int16
}

// WifiData packet returned by the Tello
type WifiData struct {
	Disturb  int16
	Strength int16
}

// Driver represents the DJI Tello drone
type Driver struct {
	name                     string
	reqAddr                  string
	reqConn                  *net.UDPConn // UDP connection to send/receive drone commands
	videoConn                *net.UDPConn // UDP connection for drone video
	respPort                 string
	cmdMutex                 sync.Mutex
	rx, ry, lx, ly, throttle float32
	gobot.Eventer
}

// NewDriver creates a driver for the Tello drone. Pass in the UDP port to use for the responses
// from the drone.
func NewDriver(port string) *Driver {
	d := &Driver{name: gobot.DefaultName("Tello"),
		reqAddr:  "192.168.10.1:8889",
		respPort: port,
		Eventer:  gobot.NewEventer(),
	}

	d.AddEvent(ConnectedEvent)
	d.AddEvent(FlightDataEvent)
	d.AddEvent(WifiDataEvent)
	d.AddEvent(LightStrengthEvent)
	d.AddEvent(VideoFrameEvent)
	d.AddEvent(SetExposureEvent)

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
	d.reqConn, err = net.DialUDP("udp", respPort, reqAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// handle responses
	go func() {
		for {
			err := d.handleResponse()
			if err != nil {
				fmt.Println("response parse error:", err)
			}
		}
	}()

	// starts notifications coming from drone to port 6038 aka 0x9617 when encoded low-endian.
	// TODO: allow setting a specific video port.
	d.SendCommand("conn_req:\x96\x17")

	// send stick commands
	go func() {
		time.Sleep(50 * time.Millisecond)
		for {
			err := d.SendStickCommand()
			if err != nil {
				fmt.Println("stick command error:", err)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	return nil
}

func (d *Driver) handleResponse() error {
	var buf [2048]byte
	n, err := d.reqConn.Read(buf[0:])
	if err != nil {
		return err
	}

	// parse binary packet
	if buf[0] == messageStart {
		if buf[6] == 0x10 {
			switch buf[5] {
			case logMessage:
				d.Publish(d.Event(LogEvent), buf[9:])
			default:
				fmt.Printf("Unknown message: %+v\n", buf[0:n])
			}
			return nil
		}

		switch buf[5] {
		case wifiMessage:
			buf := bytes.NewReader(buf[9:12])
			wd := &WifiData{}

			err = binary.Read(buf, binary.LittleEndian, &wd.Disturb)
			err = binary.Read(buf, binary.LittleEndian, &wd.Strength)
			d.Publish(d.Event(WifiDataEvent), wd)
		case lightMessage:
			buf := bytes.NewReader(buf[9:10])
			var ld int16

			err = binary.Read(buf, binary.LittleEndian, &ld)
			d.Publish(d.Event(LightStrengthEvent), ld)
		case timeMessage:
			d.Publish(d.Event(TimeEvent), buf[7:8])
		case takeoffCommand:
			d.Publish(d.Event(TakeoffEvent), buf[7:8])
		case landCommand:
			d.Publish(d.Event(LandingEvent), buf[7:8])
		case flipCommand:
			d.Publish(d.Event(FlipEvent), buf[7:8])
		case flightMessage:
			fd, _ := d.ParseFlightData(buf[9:])
			d.Publish(d.Event(FlightDataEvent), fd)
		case exposureCommand:
			d.Publish(d.Event(SetExposureEvent), buf[7:8])
		default:
			fmt.Printf("Unknown message: %+v\n", buf[0:n])
		}
		return nil
	}

	// parse text packet
	if buf[0] == 0x63 && buf[1] == 0x6f && buf[2] == 0x6e {
		d.Publish(d.Event(ConnectedEvent), nil)
		d.processVideo()
	}

	return nil
}

func (d *Driver) processVideo() error {
	videoPort, err := net.ResolveUDPAddr("udp", ":6038")
	if err != nil {
		return err
	}
	d.videoConn, err = net.ListenUDP("udp", videoPort)
	if err != nil {
		return err
	}

	go func() {
		buf := make([]byte, 2048)
		for {
			n, _, err := d.videoConn.ReadFromUDP(buf)
			d.Publish(d.Event(VideoFrameEvent), buf[2:n])

			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}()

	return nil
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.reqConn.Close()
	d.videoConn.Close()
	return
}

// TakeOff tells drones to liftoff and start flying.
func (d *Driver) TakeOff() (err error) {
	takeOffPacket := []byte{messageStart, 0x58, 0x00, 0x7c, 0x68, takeoffCommand, 0x00, 0xe4, 0x01, 0xc2, 0x16}
	_, err = d.reqConn.Write(takeOffPacket)
	return
}

// Land tells drone to come in for landing.
func (d *Driver) Land() (err error) {
	landPacket := []byte{messageStart, 0x60, 0x00, 0x27, 0x68, landCommand, 0x00, 0xe5, 0x01, 0x00, 0xba, 0xc7}
	_, err = d.reqConn.Write(landPacket)
	return
}

// StartVideo tells to start video stream.
func (d *Driver) StartVideo() (err error) {
	pkt := []byte{messageStart, 0x58, 0x00, 0x7c, 0x60, videoStartCommand, 0x00, 0x00, 0x00, 0x6c, 0x95}
	_, err = d.reqConn.Write(pkt)
	return
}

// SetExposure sets the drone camera exposure level. Valid levels are 0, 1, and 2.
func (d *Driver) SetExposure(level int) (err error) {
	if level < 0 || level > 2 {
		return errors.New("Invalid exposure level")
	}
	pkt := []byte{messageStart, 0x60, 0x00, 0x27, 0x48, exposureCommand, 0x00, 0xe6, 0x01, byte(level), 0x00, 0x00}

	// sets ending crc bytes for packet
	l := len(pkt)
	pkt[(l - 2)], pkt[(l - 1)] = CalculateCRC(pkt)

	_, err = d.reqConn.Write(pkt)
	return
}

// SetVideoEncoderRate sets the drone video encoder rate.
func (d *Driver) SetVideoEncoderRate(rate int) (err error) {
	pkt := []byte{messageStart, 0x62, 0x00, 0x27, 0x68, videoEncoderRateCommand, 0x00, 0xe6, 0x01, byte(rate), 0x00, 0x00, 0x00}

	// sets ending crc bytes for packet
	l := len(pkt)
	pkt[(l - 2)], pkt[(l - 1)] = CalculateCRC(pkt)

	_, err = d.reqConn.Write(pkt)
	return
}

// Rate does some still unknown thing.
func (d *Driver) Rate() (err error) {
	pkt := []byte{messageStart, 0x58, 0x00, 0x7c, 0x48, 40, 0x00, 0xe6, 0x01, 0x6c, 0x95}

	// sets ending crc bytes for packet
	l := len(pkt)
	pkt[(l - 2)], pkt[(l - 1)] = CalculateCRC(pkt)

	_, err = d.reqConn.Write(pkt)
	return
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

// Flip tells drone to flip
func (d *Driver) Flip(direction int) (err error) {
	pkt := []byte{messageStart, 0x60, 0x00, 0x27, 0x70, flipCommand, 0x00, 0xe6, 0x01, byte(direction), 0x00, 0x00}

	// sets ending crc bytes for packet
	l := len(pkt)
	pkt[(l - 2)], pkt[(l - 1)] = CalculateCRC(pkt)

	_, err = d.reqConn.Write(pkt)
	return
}

// FrontFlip tells the drone to perform a front flip.
func (d *Driver) FrontFlip() (err error) {
	return d.Flip(flipFront)
}

// BackFlip tells the drone to perform a back flip.
func (d *Driver) BackFlip() (err error) {
	return d.Flip(flipBack)
}

// RightFlip tells the drone to perform a flip to the right.
func (d *Driver) RightFlip() (err error) {
	return d.Flip(flipRight)
}

// LeftFlip tells the drone to perform a flip to the left.
func (d *Driver) LeftFlip() (err error) {
	return d.Flip(flipLeft)
}

// ParseFlightData from drone
func (d *Driver) ParseFlightData(b []byte) (fd *FlightData, err error) {
	buf := bytes.NewReader(b)
	fd = &FlightData{}
	var data byte

	if buf.Len() < 24 {
		err = errors.New("Invalid buffer length for flight data packet")
		fmt.Println(err)
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fd.height)
	err = binary.Read(buf, binary.LittleEndian, &fd.northSpeed)
	err = binary.Read(buf, binary.LittleEndian, &fd.eastSpeed)
	err = binary.Read(buf, binary.LittleEndian, &fd.groundSpeed)
	err = binary.Read(buf, binary.LittleEndian, &fd.flyTime)

	err = binary.Read(buf, binary.LittleEndian, &data)
	fd.imuState = int16(data >> 0 & 0x1)
	fd.pressureState = int16(data >> 1 & 0x1)
	fd.downVisualState = int16(data >> 2 & 0x1)
	fd.powerState = int16(data >> 3 & 0x1)
	fd.batteryState = int16(data >> 4 & 0x1)
	fd.gravityState = int16(data >> 5 & 0x1)
	fd.windState = int16(data >> 7 & 0x1)

	err = binary.Read(buf, binary.LittleEndian, &fd.imuCalibrationState)
	err = binary.Read(buf, binary.LittleEndian, &fd.batteryPercentage)
	err = binary.Read(buf, binary.LittleEndian, &fd.droneFlyTimeLeft)
	err = binary.Read(buf, binary.LittleEndian, &fd.droneBatteryLeft)

	err = binary.Read(buf, binary.LittleEndian, &data)
	fd.eMSky = int16(data >> 0 & 0x1)
	fd.eMGround = int16(data >> 1 & 0x1)
	fd.eMOpen = int16(data >> 2 & 0x1)
	fd.droneHover = int16(data >> 3 & 0x1)
	fd.outageRecording = int16(data >> 4 & 0x1)
	fd.batteryLow = int16(data >> 5 & 0x1)
	fd.batteryLower = int16(data >> 6 & 0x1)
	fd.factoryMode = int16(data >> 7 & 0x1)

	err = binary.Read(buf, binary.LittleEndian, &fd.flyMode)
	err = binary.Read(buf, binary.LittleEndian, &fd.throwFlyTimer)
	err = binary.Read(buf, binary.LittleEndian, &fd.cameraState)

	err = binary.Read(buf, binary.LittleEndian, &data)
	fd.electricalMachineryState = int16(data & 0xff)

	err = binary.Read(buf, binary.LittleEndian, &data)
	fd.frontIn = int16(data >> 0 & 0x1)
	fd.frontOut = int16(data >> 1 & 0x1)
	fd.frontLSC = int16(data >> 2 & 0x1)

	err = binary.Read(buf, binary.LittleEndian, &data)
	fd.temperatureHeight = int16(data >> 0 & 0x1)

	return
}

// SendStickCommand sends the joystick command packet to the drone.
func (d *Driver) SendStickCommand() (err error) {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	pkt := []byte{messageStart, 0xb0, 0x00, 0x7f, 0x60, stickCommand, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x16, 0x01, 0x0e, 0x00, 0x25, 0x54}

	// RightX center=1024 left =364 right =-364
	axis1 := int16(660.0*d.rx + 1024.0)

	//RightY down =364 up =-364
	axis2 := int16(660.0*d.ry + 1024.0)

	//LeftY down =364 up =-364
	axis3 := int16(660.0*d.ly + 1024.0)

	//LeftX left =364 right =-364
	axis4 := int16(660.0*d.lx + 1024.0)

	// speed control
	axis5 := int16(660.0*d.throttle + 1024.0)

	packedAxis := int64(axis1)&0x7FF | int64(axis2&0x7FF)<<11 | 0x7FF&int64(axis3)<<22 | 0x7FF&int64(axis4)<<33 | int64(axis5)<<44
	pkt[9] = byte(0xFF & packedAxis)
	pkt[10] = byte(packedAxis >> 8 & 0xFF)
	pkt[11] = byte(packedAxis >> 16 & 0xFF)
	pkt[12] = byte(packedAxis >> 24 & 0xFF)
	pkt[13] = byte(packedAxis >> 32 & 0xFF)
	pkt[14] = byte(packedAxis >> 40 & 0xFF)

	now := time.Now()
	pkt[15] = byte(now.Hour())
	pkt[16] = byte(now.Minute())
	pkt[17] = byte(now.Second())
	pkt[18] = byte(now.UnixNano() / int64(time.Millisecond) & 0xff)
	pkt[19] = byte(now.UnixNano() / int64(time.Millisecond) >> 8)

	// sets ending crc for packet
	l := len(pkt)
	pkt[(l - 2)], pkt[(l - 1)] = CalculateCRC(pkt)

	_, err = d.reqConn.Write(pkt)
	return
}

// SendCommand is used to send a text command such as the initial connection request to the drone.
func (d *Driver) SendCommand(cmd string) (err error) {
	_, err = d.reqConn.Write([]byte(cmd))
	return
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}
