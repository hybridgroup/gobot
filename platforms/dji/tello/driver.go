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

	// SetVideoEncoderRateEvent event
	SetVideoEncoderRateEvent = "setvideoencoder"
)

const (
	messageStart   = 0xcc
	wifiMessage    = 26
	videoRateQuery = 40
	lightMessage   = 53
	flightMessage  = 86

	logMessage = 0x50

	videoEncoderRateCommand = 0x20
	videoStartCommand       = 0x25
	exposureCommand         = 0x34
	timeCommand             = 70
	stickCommand            = 80
	takeoffCommand          = 0x54
	landCommand             = 0x55
	flipCommand             = 0x5c
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
	seq                      int16
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
	d.AddEvent(TakeoffEvent)
	d.AddEvent(LandingEvent)
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
		case timeCommand:
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
		d.SendDateTime()
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
		for {
			buf := make([]byte, 2048)
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
	buf, _ := d.createPacket(takeoffCommand, 0x68, 0)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// Land tells drone to come in for landing.
func (d *Driver) Land() (err error) {
	buf, _ := d.createPacket(landCommand, 0x68, 1)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, byte(0x00))
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// StartVideo tells Tello to send start info (SPS/PPS) for video stream.
func (d *Driver) StartVideo() (err error) {
	buf, _ := d.createPacket(videoStartCommand, 0x60, 0)
	binary.Write(buf, binary.LittleEndian, int16(0x00)) // seq = 0
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// SetExposure sets the drone camera exposure level. Valid levels are 0, 1, and 2.
func (d *Driver) SetExposure(level int) (err error) {
	if level < 0 || level > 2 {
		return errors.New("Invalid exposure level")
	}

	buf, _ := d.createPacket(exposureCommand, 0x48, 1)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, byte(level))
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// SetVideoEncoderRate sets the drone video encoder rate.
func (d *Driver) SetVideoEncoderRate(rate VideoBitRate) (err error) {
	buf, _ := d.createPacket(videoEncoderRateCommand, 0x68, 1)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, byte(rate))
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// Rate queries the current video bit rate.
func (d *Driver) Rate() (err error) {
	buf, _ := d.createPacket(videoRateQuery, 0x48, 0)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
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
func (d *Driver) Flip(direction FlipType) (err error) {
	buf, _ := d.createPacket(flipCommand, 0x70, 1)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)
	binary.Write(buf, binary.LittleEndian, byte(direction))
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// FrontFlip tells the drone to perform a front flip.
func (d *Driver) FrontFlip() (err error) {
	return d.Flip(FlipFront)
}

// BackFlip tells the drone to perform a back flip.
func (d *Driver) BackFlip() (err error) {
	return d.Flip(FlipBack)
}

// RightFlip tells the drone to perform a flip to the right.
func (d *Driver) RightFlip() (err error) {
	return d.Flip(FlipRight)
}

// LeftFlip tells the drone to perform a flip to the left.
func (d *Driver) LeftFlip() (err error) {
	return d.Flip(FlipLeft)
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

	buf, _ := d.createPacket(stickCommand, 0x60, 11)
	binary.Write(buf, binary.LittleEndian, int16(0x00)) // seq = 0

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
	binary.Write(buf, binary.LittleEndian, byte(0xFF&packedAxis))
	binary.Write(buf, binary.LittleEndian, byte(packedAxis>>8&0xFF))
	binary.Write(buf, binary.LittleEndian, byte(packedAxis>>16&0xFF))
	binary.Write(buf, binary.LittleEndian, byte(packedAxis>>24&0xFF))
	binary.Write(buf, binary.LittleEndian, byte(packedAxis>>32&0xFF))
	binary.Write(buf, binary.LittleEndian, byte(packedAxis>>40&0xFF))

	now := time.Now()
	binary.Write(buf, binary.LittleEndian, byte(now.Hour()))
	binary.Write(buf, binary.LittleEndian, byte(now.Minute()))
	binary.Write(buf, binary.LittleEndian, byte(now.Second()))
	binary.Write(buf, binary.LittleEndian, byte(now.UnixNano()/int64(time.Millisecond)&0xff))
	binary.Write(buf, binary.LittleEndian, byte(now.UnixNano()/int64(time.Millisecond)>>8))

	// sets ending crc for packet
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())

	return
}

// SendDateTime sends the current date/time to the drone.
func (d *Driver) SendDateTime() (err error) {
	d.cmdMutex.Lock()
	defer d.cmdMutex.Unlock()

	buf, _ := d.createPacket(timeCommand, 0x50, 11)
	d.seq++
	binary.Write(buf, binary.LittleEndian, d.seq)

	now := time.Now()
	binary.Write(buf, binary.LittleEndian, byte(0x00))
	binary.Write(buf, binary.LittleEndian, now.Hour())
	binary.Write(buf, binary.LittleEndian, now.Minute())
	binary.Write(buf, binary.LittleEndian, now.Second())
	binary.Write(buf, binary.LittleEndian, int16(now.UnixNano()/int64(time.Millisecond)&0xff))
	binary.Write(buf, binary.LittleEndian, int16(now.UnixNano()/int64(time.Millisecond)>>8))

	// sets ending crc for packet
	binary.Write(buf, binary.LittleEndian, CalculateCRC16(buf.Bytes()))

	_, err = d.reqConn.Write(buf.Bytes())
	return
}

// SendCommand is used to send a text command such as the initial connection request to the drone.
func (d *Driver) SendCommand(cmd string) (err error) {
	_, err = d.reqConn.Write([]byte(cmd))
	return
}

func (d *Driver) createPacket(cmd int16, pktType byte, len int16) (buf *bytes.Buffer, err error) {
	l := len + 11
	buf = &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, byte(messageStart))
	binary.Write(buf, binary.LittleEndian, l<<3)
	binary.Write(buf, binary.LittleEndian, CalculateCRC8(buf.Bytes()[0:3]))
	binary.Write(buf, binary.LittleEndian, pktType)
	binary.Write(buf, binary.LittleEndian, cmd)

	return buf, nil
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	}
	if val < 0 {
		return 0
	}

	return val
}
