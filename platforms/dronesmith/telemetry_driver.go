package dronesmith

import (
	"time"

	"github.com/hybridgroup/gobot"
)

type TelemetryDriver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Commander
}

func NewTelemetryDriver(a *Adaptor) *TelemetryDriver {
	d := &TelemetryDriver{
		name:       "Telemetry",
		connection: a,
		interval:   500 * time.Millisecond,
		halt:       make(chan bool, 0),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *TelemetryDriver) Name() string { return d.name }

func (d *TelemetryDriver) SetName(n string) { d.name = n }

func (d *TelemetryDriver) Connection() gobot.Connection {
	return d.connection
}

func (d *TelemetryDriver) adaptor() *Adaptor {
	return d.Connection().(*Adaptor)
}

func (d *TelemetryDriver) Start() (err error) {
	return
}

func (d *TelemetryDriver) Halt() (err error) {
	return
}

// Info reads general Drone info from Dronesmith cloud api
func (d *TelemetryDriver) Info() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/info", nil)
	return
}

// Status gets the current Drone status from Dronesmith cloud api
func (d *TelemetryDriver) Status() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/status", nil)
	return
}

// Mode gets the current Drone mode from Dronesmith cloud api
func (d *TelemetryDriver) Mode() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/mode", nil)
	return
}

// GPS gets the current Drone GPS position from Dronesmith cloud api
func (d *TelemetryDriver) GPS() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/gps", nil)
	return
}

// Attitude gets the current Drone attitude info from Dronesmith cloud api
func (d *TelemetryDriver) Attitude() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/attitude", nil)
	return
}

// Position gets the current Drone position info from Dronesmith cloud api
func (d *TelemetryDriver) Position() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/position", nil)
	return
}

// Sensors gets the current Drone Sensors info from Dronesmith cloud api
func (d *TelemetryDriver) Sensors() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/sensors", nil)
	return
}

// Home gets the current Drone Home info from Dronesmith cloud api
func (d *TelemetryDriver) Home() (m map[string]interface{}, err error) {
	m, err = d.adaptor().Request("GET", "/home", nil)
	return
}
