package tello

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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
			err := d.handleResponse()
			if err != nil {
				return
			}
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

func (d *Driver) handleResponse() error {
	var buf [256]byte
	n, err := d.reqConn.Read(buf[0:])
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp := string(buf[0:n])
	d.responses <- resp
	return nil
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
		switch res {
		case "OK":
			return nil
		case "FALSE":
			return errors.New("Command returned false: " + cmd)
		default:
			fmt.Println("Unknown response:", res)
			return nil
		}
	case <-time.After(5 * time.Second):
		return errors.New("Command timeout: " + cmd)
	}
}

func (d *Driver) sendFunction(cmd string) (string, error) {
	_, err := d.reqConn.Write([]byte(cmd))
	if err != nil {
		return "", err
	}

	select {
	case res := <-d.responses:
		return strings.Replace(res, "\r\n", "", -1), nil
	case <-time.After(5 * time.Second):
		return "", errors.New("Command timeout: " + cmd)
	}
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

// FrontFlip tells the drone to perform a front flip
func (d *Driver) FrontFlip() (err error) {
	return d.sendCommand("flip f")
}

// BackFlip tells the drone to perform a backflip
func (d *Driver) BackFlip() (err error) {
	return d.sendCommand("flip b")
}

// RightFlip tells the drone to perform a flip to the right
func (d *Driver) RightFlip() (err error) {
	return d.sendCommand("flip r")
}

// LeftFlip tells the drone to perform a flip to the left
func (d *Driver) LeftFlip() (err error) {
	return d.sendCommand("flip l")
}

// StartRecording is not yet supported.
func (d *Driver) StartRecording() error {
	return nil
}

// StopRecording is not yet supported.
func (d *Driver) StopRecording() error {
	return nil
}

// Battery returns the current battery level in the drone as a percentage.
func (d *Driver) Battery() (int, error) {
	res, err := d.sendFunction("battery?")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// FlightTime returns the current elapsed flight time for the drone in seconds.
func (d *Driver) FlightTime() (string, error) {
	return d.sendFunction("time?")
}

// Speed returns the current speed for the drone in cm per second.
func (d *Driver) Speed() (float64, error) {
	res, err := d.sendFunction("speed?")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(res, 32)
}

// SetSpeed sets the drone speed from 1-100 cm per second.
func (d *Driver) SetSpeed(speed int) error {
	return d.sendCommand(fmt.Sprintf("speed %d", speed))
}
