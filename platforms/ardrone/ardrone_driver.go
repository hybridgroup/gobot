package ardrone

import (
	"github.com/hybridgroup/gobot"
)

type ArdroneDriver struct {
	gobot.Driver
	Adaptor DroneInterface
}

type DroneInterface interface {
	Drone() drone
}

func NewArdroneDriver(adaptor DroneInterface, name string) *ArdroneDriver {
	return &ArdroneDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]chan interface{}{
				"Flying": make(chan interface{}, 1),
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
func (a *ArdroneDriver) Init() bool {
	return true
}

func (a *ArdroneDriver) TakeOff() {
	gobot.Publish(a.Events["Flying"], gobot.Call(a.Adaptor.Drone(), "Takeoff"))
}
func (a *ArdroneDriver) Land() {
	gobot.Call(a.Adaptor.Drone(), "Land")
}
func (a *ArdroneDriver) Up(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Up", n)
}
func (a *ArdroneDriver) Down(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Down", n)
}
func (a *ArdroneDriver) Left(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Left", n)
}
func (a *ArdroneDriver) Right(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Right", n)
}
func (a *ArdroneDriver) Forward(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Forward", n)
}
func (a *ArdroneDriver) Backward(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Backward", n)
}
func (a *ArdroneDriver) Clockwise(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Clockwise", n)
}
func (a *ArdroneDriver) CounterClockwise(n float64) {
	gobot.Call(a.Adaptor.Drone(), "Counterclockwise", n)
}
func (a *ArdroneDriver) Hover() {
	gobot.Call(a.Adaptor.Drone(), "Hover")
}
