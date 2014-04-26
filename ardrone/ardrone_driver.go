package gobotArdrone

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

func NewArdrone(adaptor DroneInterface) *ArdroneDriver {
	d := new(ArdroneDriver)
	d.Adaptor = adaptor
	d.Events = make(map[string]chan interface{})
	d.Events["Flying"] = make(chan interface{}, 1)
	d.Commands = []string{}
	return d
}

func (me *ArdroneDriver) Start() bool {
	return true
}
func (me *ArdroneDriver) Halt() bool {
	return true
}
func (me *ArdroneDriver) Init() bool {
	return true
}

func (me *ArdroneDriver) TakeOff() {
	gobot.Publish(me.Events["Flying"], gobot.Call(me.Adaptor.Drone(), "Takeoff"))
}
func (me *ArdroneDriver) Land() {
	gobot.Call(me.Adaptor.Drone(), "Land")
}
func (me *ArdroneDriver) Up(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Up", a)
}
func (me *ArdroneDriver) Down(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Down", a)
}
func (me *ArdroneDriver) Left(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Left", a)
}
func (me *ArdroneDriver) Right(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Right", a)
}
func (me *ArdroneDriver) Forward(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Forward", a)
}
func (me *ArdroneDriver) Backward(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Backward", a)
}
func (me *ArdroneDriver) Clockwise(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Clockwise", a)
}
func (me *ArdroneDriver) CounterClockwise(a float64) {
	gobot.Call(me.Adaptor.Drone(), "Counterclockwise", a)
}
func (me *ArdroneDriver) Hover() {
	gobot.Call(me.Adaptor.Drone(), "Hover")
}
