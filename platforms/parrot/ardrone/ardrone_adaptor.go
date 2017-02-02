package ardrone

import (
	client "github.com/hybridgroup/go-ardrone/client"
	"gobot.io/x/gobot"
)

// drone defines expected drone behaviour
type drone interface {
	Takeoff() bool
	Land()
	Up(n float64)
	Down(n float64)
	Left(n float64)
	Right(n float64)
	Forward(n float64)
	Backward(n float64)
	Clockwise(n float64)
	Counterclockwise(n float64)
	Hover()
}

// Adaptor is gobot.Adaptor representation for the Ardrone
type Adaptor struct {
	name    string
	drone   drone
	config  client.Config
	connect func(*Adaptor) (drone, error)
}

// NewAdaptor returns a new ardrone.Adaptor and optionally accepts:
//
//  string: The ardrones IP Address
//
func NewAdaptor(v ...string) *Adaptor {
	a := &Adaptor{
		name: gobot.DefaultName("ARDrone"),
		connect: func(a *Adaptor) (drone, error) {
			return client.Connect(a.config)
		},
	}

	a.config = client.DefaultConfig()
	if len(v) > 0 {
		a.config.Ip = v[0]
	}

	return a
}

// Name returns the Adaptor Name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Adaptor Name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect establishes a connection to the ardrone
func (a *Adaptor) Connect() (err error) {
	d, err := a.connect(a)
	if err != nil {
		return err
	}
	a.drone = d
	return
}

// Finalize terminates the connection to the ardrone
func (a *Adaptor) Finalize() (err error) { return }
