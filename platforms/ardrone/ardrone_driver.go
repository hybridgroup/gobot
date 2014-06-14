package ardrone

import (
	"github.com/hybridgroup/gobot"
)

type ArdroneDriver struct {
	gobot.Driver
	Adaptor *ArdroneAdaptor
}

func NewArdroneDriver(adaptor *ArdroneAdaptor, name string) *ArdroneDriver {
	return &ArdroneDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]*gobot.Event{
				"Flying": gobot.NewEvent(),
			},
		},
		Adaptor: adaptor,
	}
}

func (a *ArdroneDriver) Start() bool {
	return true
}

func (a *ArdroneDriver) Halt() bool {
	return true
}

func (a *ArdroneDriver) TakeOff() {
	gobot.Publish(a.Events["Flying"], a.Adaptor.drone.Takeoff())
}

func (a *ArdroneDriver) Land() {
	a.Adaptor.drone.Land()
}

func (a *ArdroneDriver) Up(n float64) {
	a.Adaptor.drone.Up(n)
}

func (a *ArdroneDriver) Down(n float64) {
	a.Adaptor.drone.Down(n)
}

func (a *ArdroneDriver) Left(n float64) {
	a.Adaptor.drone.Left(n)
}

func (a *ArdroneDriver) Right(n float64) {
	a.Adaptor.drone.Right(n)
}

func (a *ArdroneDriver) Forward(n float64) {
	a.Adaptor.drone.Forward(n)
}

func (a *ArdroneDriver) Backward(n float64) {
	a.Adaptor.drone.Backward(n)
}

func (a *ArdroneDriver) Clockwise(n float64) {
	a.Adaptor.drone.Clockwise(n)
}

func (a *ArdroneDriver) CounterClockwise(n float64) {
	a.Adaptor.drone.Counterclockwise(n)
}

func (a *ArdroneDriver) Hover() {
	a.Adaptor.drone.Hover()
}
