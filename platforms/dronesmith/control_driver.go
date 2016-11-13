package dronesmith

import (
	"time"

	"github.com/hybridgroup/gobot"
)

type ControlDriver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Commander
}

func NewControlDriver(a *Adaptor) *ControlDriver {
	d := &ControlDriver{
		name:       "Control",
		connection: a,
		interval:   500 * time.Millisecond,
		halt:       make(chan bool, 0),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *ControlDriver) Name() string { return d.name }

func (d *ControlDriver) SetName(n string) { d.name = n }

func (d *ControlDriver) Connection() gobot.Connection {
	return d.connection
}

func (d *ControlDriver) adaptor() *Adaptor {
	return d.Connection().(*Adaptor)
}

func (d *ControlDriver) Start() (err error) {
	return
}

func (d *ControlDriver) Halt() (err error) {
	return
}

// Arm tells Drone to get ready to fly
func (d *ControlDriver) Arm() (err error) {
	_, err = d.adaptor().Request("POST", "/arm", nil)
	return
}

// Disarm tells Drone to prevent flight
func (d *ControlDriver) Disarm() (err error) {
	_, err = d.adaptor().Request("POST", "/disarm", nil)
	return
}

// Takeoff tells Drone to takeoff
func (d *ControlDriver) Takeoff() (err error) {
	_, err = d.adaptor().Request("POST", "/takeoff", nil)
	return
}

// Land tells Drone to land
func (d *ControlDriver) Land() (err error) {
	_, err = d.adaptor().Request("POST", "/land", nil)
	return
}
