package ardrone

import (
	"gobot.io/x/gobot"
)

const (
	// Flying event
	Flying = "flying"
)

// Driver is gobot.Driver representation for the Ardrone
type Driver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewDriver creates an Driver for the ARDrone.
//
// It add the following events:
//     'flying' - Sent when the device has taken off.
func NewDriver(connection *Adaptor) *Driver {
	d := &Driver{
		name:       gobot.DefaultName("ARDrone"),
		connection: connection,
		Eventer:    gobot.NewEventer(),
	}
	d.AddEvent(Flying)
	return d
}

// Name returns the Driver Name
func (a *Driver) Name() string { return a.name }

// SetName sets the Driver Name
func (a *Driver) SetName(n string) { a.name = n }

// Connection returns the Driver Connection
func (a *Driver) Connection() gobot.Connection { return a.connection }

// adaptor returns ardrone adaptor
func (a *Driver) adaptor() *Adaptor {
	return a.Connection().(*Adaptor)
}

// Start starts the Driver
func (a *Driver) Start() (err error) {
	return
}

// Halt halts the Driver
func (a *Driver) Halt() (err error) {
	return
}

// TakeOff makes the drone start flying, and publishes `flying` event
func (a *Driver) TakeOff() {
	a.Publish(a.Event("flying"), a.adaptor().drone.Takeoff())
}

// Land causes the drone to land
func (a *Driver) Land() {
	a.adaptor().drone.Land()
}

// Up makes the drone gain altitude.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Up(speed float64) {
	a.adaptor().drone.Up(speed)
}

// Down makes the drone reduce altitude.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Down(speed float64) {
	a.adaptor().drone.Down(speed)
}

// Left causes the drone to bank to the left, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Left(speed float64) {
	a.adaptor().drone.Left(speed)
}

// Right causes the drone to bank to the right, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Right(speed float64) {
	a.adaptor().drone.Right(speed)
}

// Forward causes the drone go forward, controls the pitch.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Forward(speed float64) {
	a.adaptor().drone.Forward(speed)
}

// Backward causes the drone go backward, controls the pitch.
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Backward(speed float64) {
	a.adaptor().drone.Backward(speed)
}

// Clockwise causes the drone to spin in clockwise direction
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) Clockwise(speed float64) {
	a.adaptor().drone.Clockwise(speed)
}

// CounterClockwise the drone to spin in counter clockwise direction
// speed can be a value from `0.0` to `1.0`.
func (a *Driver) CounterClockwise(speed float64) {
	a.adaptor().drone.Counterclockwise(speed)
}

// Hover makes the drone to hover in place.
func (a *Driver) Hover() {
	a.adaptor().drone.Hover()
}
