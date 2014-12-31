package ardrone

import (
	client "github.com/hybridgroup/go-ardrone/client"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*ArdroneAdaptor)(nil)

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

// ArdroneAdaptor is gobot.Adaptor representation for the Ardrone
type ArdroneAdaptor struct {
	name    string
	drone   drone
	config  client.Config
	connect func(*ArdroneAdaptor) (drone, error)
}

// NewArdroneAdaptor returns a new ArdroneAdaptor and optionally accepts:
//
//  string: The ardrones IP Address
//
func NewArdroneAdaptor(name string, v ...string) *ArdroneAdaptor {
	a := &ArdroneAdaptor{
		name: name,
		connect: func(a *ArdroneAdaptor) (drone, error) {
			return client.Connect(a.config)
		},
	}

	a.config = client.DefaultConfig()
	if len(v) > 0 {
		a.config.Ip = v[0]
	}

	return a
}

// Name returns the ArdroneAdaptors Name
func (a *ArdroneAdaptor) Name() string { return a.name }

// Connect establishes a connection to the ardrone
func (a *ArdroneAdaptor) Connect() (errs []error) {
	d, err := a.connect(a)
	if err != nil {
		return []error{err}
	}
	a.drone = d
	return
}

// Finalize terminates the connection to the ardrone
func (a *ArdroneAdaptor) Finalize() (errs []error) { return }
