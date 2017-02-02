package bebop

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/bebop/client"
)

// drone defines expected drone behaviour
type drone interface {
	TakeOff() error
	Land() error
	Up(n int) error
	Down(n int) error
	Left(n int) error
	Right(n int) error
	Forward(n int) error
	Backward(n int) error
	Clockwise(n int) error
	CounterClockwise(n int) error
	Stop() error
	Connect() error
	Video() chan []byte
	StartRecording() error
	StopRecording() error
	HullProtection(protect bool) error
	Outdoor(outdoor bool) error
	VideoEnable(enable bool) error
	VideoStreamMode(mode int8) error
}

// Adaptor is gobot.Adaptor representation for the Bebop
type Adaptor struct {
	name    string
	drone   drone
	connect func(*Adaptor) error
}

// NewAdaptor returns a new BebopAdaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:  gobot.DefaultName("Bebop"),
		drone: client.New(),
		connect: func(a *Adaptor) error {
			return a.drone.Connect()
		},
	}
}

// Name returns the Bebop Adaptors Name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Bebop Adaptors Name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect establishes a connection to the ardrone
func (a *Adaptor) Connect() (err error) {
	err = a.connect(a)
	return
}

// Finalize terminates the connection to the ardrone
func (a *Adaptor) Finalize() (err error) { return }
