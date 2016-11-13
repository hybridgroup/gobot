package dronesmith

import (
	"time"

	"github.com/hybridgroup/gobot"
)

type Driver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Commander
}

func NewTelemetryDriver(a *Adaptor) *Driver {
	d := &Driver{
		name:       "Telemetry",
		connection: a,
		interval:   500 * time.Millisecond,
		halt:       make(chan bool, 0),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *Driver) Name() string { return d.name }

func (d *Driver) SetName(n string) { d.name = n }

func (d *Driver) Connection() gobot.Connection {
	return d.connection
}

func (d *Driver) adaptor() *Adaptor {
	return d.Connection().(*Adaptor)
}

func (d *Driver) Start() (err error) {
	return
}

func (d *Driver) Halt() (err error) {
	return
}

// Info reads general Drone info from Dronesmith cloud api
func (d *Driver) Info() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/info", nil)
	return
}

// Status gets the current Drone status from Dronesmith cloud api
func (d *Driver) Status() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/status", nil)
	return
}

// Mode gets the current Drone mode from Dronesmith cloud api
func (d *Driver) Mode() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/mode", nil)
	return
}

// GPS gets the current Drone GPS position from Dronesmith cloud api
func (d *Driver) GPS() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/gps", nil)
	return
}

// Attitude gets the current Drone attitude info from Dronesmith cloud api
func (d *Driver) Attitude() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/attitude", nil)
	return
}

// Position gets the current Drone position info from Dronesmith cloud api
func (d *Driver) Position() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/position", nil)
	return
}

// Sensors gets the current Drone Sensors info from Dronesmith cloud api
func (d *Driver) Sensors() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/sensors", nil)
	return
}

// Home gets the current Drone Home info from Dronesmith cloud api
func (d *Driver) Home() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/home", nil)
	return
}
