package bebop

import (
	"gobot.io/x/gobot"
)

const (
	// Flying event
	Flying = "flying"
)

// Driver is gobot.Driver representation for the Bebop
type Driver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewDriver creates an Bebop Driver.
func NewDriver(connection *Adaptor) *Driver {
	d := &Driver{
		name:       gobot.DefaultName("Bebop"),
		connection: connection,
		Eventer:    gobot.NewEventer(),
	}
	d.AddEvent(Flying)
	return d
}

// Name returns the Bebop Drivers Name
func (a *Driver) Name() string { return a.name }

// SetName sets the Bebop Drivers Name
func (a *Driver) SetName(n string) { a.name = n }

// Connection returns the Bebop Drivers Connection
func (a *Driver) Connection() gobot.Connection { return a.connection }

// adaptor returns ardrone adaptor
func (a *Driver) adaptor() *Adaptor {
	return a.Connection().(*Adaptor)
}

// Start starts the Bebop Driver
func (a *Driver) Start() (err error) {
	return
}

// Halt halts the Bebop Driver
func (a *Driver) Halt() (err error) {
	return
}

// TakeOff makes the drone start flying
func (a *Driver) TakeOff() {
	a.Publish(a.Event("flying"), a.adaptor().drone.TakeOff())
}

// Land causes the drone to land
func (a *Driver) Land() {
	a.adaptor().drone.Land()
}

// Up makes the drone gain altitude.
// speed can be a value from `0` to `100`.
func (a *Driver) Up(speed int) {
	a.adaptor().drone.Up(speed)
}

// Down makes the drone reduce altitude.
// speed can be a value from `0` to `100`.
func (a *Driver) Down(speed int) {
	a.adaptor().drone.Down(speed)
}

// Left causes the drone to bank to the left, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0` to `100`.
func (a *Driver) Left(speed int) {
	a.adaptor().drone.Left(speed)
}

// Right causes the drone to bank to the right, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0` to `100`.
func (a *Driver) Right(speed int) {
	a.adaptor().drone.Right(speed)
}

// Forward causes the drone go forward, controls the pitch.
// speed can be a value from `0` to `100`.
func (a *Driver) Forward(speed int) {
	a.adaptor().drone.Forward(speed)
}

// Backward causes the drone go forward, controls the pitch.
// speed can be a value from `0` to `100`.
func (a *Driver) Backward(speed int) {
	a.adaptor().drone.Backward(speed)
}

// Clockwise causes the drone to spin in clockwise direction
// speed can be a value from `0` to `100`.
func (a *Driver) Clockwise(speed int) {
	a.adaptor().drone.Clockwise(speed)
}

// CounterClockwise the drone to spin in counter clockwise direction
// speed can be a value from `0` to `100`.
func (a *Driver) CounterClockwise(speed int) {
	a.adaptor().drone.CounterClockwise(speed)
}

// Stop makes the drone to hover in place.
func (a *Driver) Stop() {
	a.adaptor().drone.Stop()
}

// Video returns a channel which raw video frames will be broadcast on
func (a *Driver) Video() chan []byte {
	return a.adaptor().drone.Video()
}

// StartRecording starts the recording video to the drones interal storage
func (a *Driver) StartRecording() error {
	return a.adaptor().drone.StartRecording()
}

// StopRecording stops a previously started recording
func (a *Driver) StopRecording() error {
	return a.adaptor().drone.StopRecording()
}

// HullProtection tells the drone if the hull/prop protectors are attached. This is needed to adjust flight characteristics of the Bebop.
func (a *Driver) HullProtection(protect bool) error {
	return a.adaptor().drone.HullProtection(protect)
}

// Outdoor tells the drone if flying Outdoor or not. This is needed to adjust flight characteristics of the Bebop.
func (a *Driver) Outdoor(outdoor bool) error {
	return a.adaptor().drone.Outdoor(outdoor)
}

// VideoEnable tells the drone to start/stop streaming video
func (a *Driver) VideoEnable(enable bool) error {
	return a.adaptor().drone.VideoEnable(enable)
}

// VideoStreamMode tells the drone what mode to use for streaming video
func (a *Driver) VideoStreamMode(mode int8) error {
	return a.adaptor().drone.VideoStreamMode(mode)
}
