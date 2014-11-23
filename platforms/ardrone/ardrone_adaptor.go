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

type ArdroneAdaptor struct {
	name    string
	drone   drone
	connect func(*ArdroneAdaptor) (err error)
}

// NewArdroneAdaptor creates a new ardrone and connects with default configuration
func NewArdroneAdaptor(name string, v ...string) *ArdroneAdaptor {
	return &ArdroneAdaptor{
		name: name,
		connect: func(a *ArdroneAdaptor) (err error) {
			config := client.DefaultConfig()
			if len(v) > 0 {
				config.Ip = v[0]
			}
			d, err := client.Connect(config)
			if err != nil {
				return
			}
			a.drone = d
			return
		},
	}
}

func (a *ArdroneAdaptor) Name() string { return a.name }

// Connect returns true when connection to ardrone is established correclty
func (a *ArdroneAdaptor) Connect() (errs []error) {
	if err := a.connect(a); err != nil {
		return []error{err}
	}
	return
}

// Finalize returns true when connection is finalized correctly
func (a *ArdroneAdaptor) Finalize() (errs []error) { return }
