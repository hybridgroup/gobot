package hs200

import (
	"net"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

// Driver represents the control information for the hs200 drone
type Driver struct {
	name       string
	mutex      *sync.RWMutex // Protect the command from concurrent access
	stopc      chan struct{} // Stop the flight loop goroutine
	cmd        []byte        // the UDP command packet we keep sending the drone
	enabled    bool          // Are we in an enabled state
	tcpaddress string
	udpaddress string
	udpconn    net.Conn // UDP connection to the drone
	tcpconn    net.Conn // TCP connection to the drone
}

// NewDriver creates a driver for the HolyStone hs200
func NewDriver(tcpaddress string, udpaddress string) *Driver {
	command := []byte{
		0xff, // 2 byte header
		0x04,

		// Left joystick
		0x7e, // throttle 0x00 - 0xff(?)
		0x3f, // rotate left/right

		// Right joystick
		0xc0, // forward / backward 0x80 - 0xfe(?)
		0x3f, // left / right 0x00 - 0x7e(?)

		// Trim
		0x90, // ? yaw (used as a setting to trim the yaw of the uav)
		0x10, // ? pitch (used as a setting to trim the pitch of the uav)
		0x10, // ? roll (used as a setting to trim the roll of the uav)

		0x00, // flags/buttons
		0x00, // checksum; 255 - ((sum of flight controls from index 1 to 9) % 256)
	}
	command[10] = checksum(command)

	return &Driver{name: gobot.DefaultName("HS200"), stopc: make(chan struct{}),
		tcpaddress: tcpaddress, udpaddress: udpaddress, cmd: command, mutex: &sync.RWMutex{}}
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection { return nil }

// Start starts the driver.
func (d *Driver) Start() (err error) {
	tc, terr := net.Dial("tcp", d.tcpaddress)
	if terr != nil {
		return terr
	}
	uc, uerr := net.Dial("udp4", d.udpaddress)
	if uerr != nil {
		return uerr
	}

	d.udpconn = uc
	d.tcpconn = tc

	return
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.stop()
	return
}

func (d *Driver) stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.enabled = false
}

func (d *Driver) flightLoop(stopc chan struct{}) {
	udpTick := time.NewTicker(50 * time.Millisecond)
	defer udpTick.Stop()
	tcpTick := time.NewTicker(1000 * time.Millisecond)
	defer tcpTick.Stop()
	for {
		select {
		case <-udpTick.C:
			d.sendUDP()
		case <-tcpTick.C:
			// Send TCP commands from here once we figure out what they do...
		case <-stopc:
			d.stop()
			return
		}
	}
}

func checksum(c []byte) byte {
	var sum byte
	for i := 1; i < 10; i++ {
		sum += c[i]
	}
	return 255 - sum
}

func (d *Driver) sendUDP() {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	d.udpconn.Write(d.cmd)
}

// Enable enables the drone to start flying.
func (d Driver) Enable() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.enabled {
		go d.flightLoop(d.stopc)
		d.enabled = true
	}
}

// Disable disables the drone from flying.
func (d Driver) Disable() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.enabled {
		d.stopc <- struct{}{}
	}
}

// TakeOff tells drones to liftoff and start flying.
func (d Driver) TakeOff() {
	d.mutex.Lock()
	d.cmd[9] = 0x40
	d.cmd[10] = checksum(d.cmd)
	d.mutex.Unlock()
	time.Sleep(500 * time.Millisecond)
	d.mutex.Lock()
	d.cmd[9] = 0x04
	d.cmd[10] = checksum(d.cmd)
	d.mutex.Unlock()
}

// Land tells drone to come in for landing.
func (d Driver) Land() {
	d.mutex.Lock()
	d.cmd[9] = 0x80
	d.cmd[10] = checksum(d.cmd)
	d.mutex.Unlock()
	time.Sleep(500 * time.Millisecond)
	d.mutex.Lock()
	d.cmd[9] = 0x04
	d.cmd[10] = checksum(d.cmd)
	d.mutex.Unlock()
}

// floatToCmdByte converts a float in the range of -1 to +1 to an integer command
func floatToCmdByte(cmd float32, mid byte, maxv byte) byte {
	if cmd > 1.0 {
		cmd = 1.0
	}
	if cmd < -1.0 {
		cmd = -1.0
	}
	cmd = cmd * float32(maxv)
	bval := byte(cmd + float32(mid) + 0.5)
	return bval
}

// Throttle sends the drone up from a hover (or down if speed is negative)
func (d *Driver) Throttle(speed float32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[2] = floatToCmdByte(speed, 0x7e, 0x7e)
	d.cmd[10] = checksum(d.cmd)
}

// Rotate rotates the drone (yaw)
func (d *Driver) Rotate(speed float32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[3] = floatToCmdByte(speed, 0x3f, 0x3f)
	d.cmd[10] = checksum(d.cmd)
}

// Forward sends the drone forward (or backwards if speed is negative, pitch the drone)
func (d *Driver) Forward(speed float32) {
	speed = -speed
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[4] = floatToCmdByte(speed, 0xc0, 0x3f)
	d.cmd[10] = checksum(d.cmd)
}

// Right moves the drone to the right (or left if speed is negative, rolls the drone)
func (d *Driver) Right(speed float32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[5] = floatToCmdByte(speed, 0x3f, 0x3f)
	d.cmd[10] = checksum(d.cmd)
}
