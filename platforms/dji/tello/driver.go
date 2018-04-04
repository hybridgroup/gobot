package tello

import (
	"errors"
	"fmt"
	"net"
	"time"

	"gobot.io/x/gobot"
)

// Driver represents the DJI Tello drone
type Driver struct {
	name      string
	reqAddr   string
	reqConn   *net.UDPConn // UDP connection to send/receive drone commands
	respPort  string
	responses chan string
}

// NewDriver creates a driver for the Tello drone. Pass in the UDP port to use for the responses
// from the drone.
func NewDriver(port string) *Driver {
	return &Driver{name: gobot.DefaultName("Tello"),
		reqAddr:   "192.168.10.1:8889",
		respPort:  port,
		responses: make(chan string)}
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

	go func() {
		for {
			d.handleResponse()
		}
	}()

	// puts Tello drone into command mode, so we can send it further commands
	err = d.sendCommand("command")
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Driver) handleResponse() {
	var buf [256]byte
	n, err := d.reqConn.Read(buf[0:])
	if err != nil {
		fmt.Println("Error on response")
		return
	}

	resp := string(buf[0:n])
	d.responses <- resp
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.reqConn.Close()
	return
}

func (d *Driver) sendCommand(cmd string) error {
	_, err := d.reqConn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	select {
	case res := <-d.responses:
		fmt.Println(res)
	case <-time.After(5 * time.Second):
		return errors.New("Command timeout: " + cmd)
	}
	return nil
}

// TakeOff tells drones to liftoff and start flying.
func (d *Driver) TakeOff() error {
	return d.sendCommand("takeoff")
}

// Land tells drone to come in for landing.
func (d *Driver) Land() error {
	return d.sendCommand("land")
}

// Move tells drone to move in particular direction for particular distance
func (d *Driver) Move(dir string, dist int) error {
	return d.sendCommand(fmt.Sprintf("%s %d", dir, dist))
}

// Forward sends the drone forward
func (d *Driver) Forward(dist int) error {
	return d.Move("forward", dist)
}

// Backward sends the drone backward
func (d *Driver) Backward(dist int) error {
	return d.Move("back", dist)
}

// Right sends the drone right.
func (d *Driver) Right(dist int) error {
	return d.Move("right", dist)
}

// Left sends the drone left.
func (d *Driver) Left(dist int) error {
	return d.Move("left", dist)
}

// Up sends the drone up.
func (d *Driver) Up(dist int) error {
	return d.Move("up", dist)
}

// Down sends the drone down.
func (d *Driver) Down(dist int) error {
	return d.Move("down", dist)
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 1-360.
func (d *Driver) Clockwise(deg int) error {
	return d.Move("cw", deg)
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 1-360.
func (d *Driver) CounterClockwise(deg int) error {
	return d.Move("ccw", deg)
}
