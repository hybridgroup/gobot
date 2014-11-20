package ardrone

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*ArdroneDriver)(nil)

type ArdroneDriver struct {
	gobot.Driver
}

// NewArdroneDriver creates an ArdroneDriver with specified name.
//
// It add the following events:
//     'flying' - Sent when the device has taken off.
func NewArdroneDriver(adaptor *ArdroneAdaptor, name string) *ArdroneDriver {
	d := &ArdroneDriver{
		Driver: *gobot.NewDriver(
			name,
			"ArdroneDriver",
			adaptor,
		),
	}
	d.AddEvent("flying")
	return d
}

// adaptor returns ardrone adaptor
func (a *ArdroneDriver) adaptor() *ArdroneAdaptor {
	return a.Adaptor().(*ArdroneAdaptor)
}

// Start returns true if driver is started succesfully
func (a *ArdroneDriver) Start() (errs []error) {
	return
}

// Halt returns true if driver is halted succesfully
func (a *ArdroneDriver) Halt() (errs []error) {
	return
}

// TakeOff makes the drone start flying
// and publishes `flying` event
func (a *ArdroneDriver) TakeOff() {
	gobot.Publish(a.Event("flying"), a.adaptor().drone.Takeoff())
}

// Land makes the drone stop flying
func (a *ArdroneDriver) Land() {
	a.adaptor().drone.Land()
}

// Up makes the drone gain altitude.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Up(speed float64) {
	a.adaptor().drone.Up(speed)
}

// Down makes the drone reduce altitude.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Down(speed float64) {
	a.adaptor().drone.Down(speed)
}

// Left causes the drone to bank to the left, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Left(speed float64) {
	a.adaptor().drone.Left(speed)
}

// Right causes the drone to bank to the right, controls the roll, which is
// a horizontal movement using the camera as a reference point.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Right(speed float64) {
	a.adaptor().drone.Right(speed)
}

// Forward causes the drone go forward, controls the pitch.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Forward(speed float64) {
	a.adaptor().drone.Forward(speed)
}

// Backward causes the drone go forward, controls the pitch.
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Backward(speed float64) {
	a.adaptor().drone.Backward(speed)
}

// Clockwise causes the drone to spin in clockwise direction
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) Clockwise(speed float64) {
	a.adaptor().drone.Clockwise(speed)
}

// CounterClockwise the drone to spin in counter clockwise direction
// speed can be a value from `0.0` to `1.0`.
func (a *ArdroneDriver) CounterClockwise(speed float64) {
	a.adaptor().drone.Counterclockwise(speed)
}

// Hover makes the drone to hover in place.
func (a *ArdroneDriver) Hover() {
	a.adaptor().drone.Hover()
}
