package ardrone

import (
	client "github.com/hybridgroup/go-ardrone/client"
	"github.com/hybridgroup/gobot"
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

type ArdroneAdaptor struct {
	gobot.Adaptor
	drone   drone
	connect func(*ArdroneAdaptor)
}

// NewArdroneAdaptor creates a new ardrone and connects with default configuration
func NewArdroneAdaptor(name string) *ArdroneAdaptor {
	return &ArdroneAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"ArdroneAdaptor",
		),
		connect: func(a *ArdroneAdaptor) {
			d, err := client.Connect(client.DefaultConfig())
			if err != nil {
				panic(err)
			}
			a.drone = d
		},
	}
}

// Connect returns true when connection to ardrone is established correclty
func (a *ArdroneAdaptor) Connect() bool {
	a.connect(a)
	return true
}

// Finalize returns true when connection is finalized correctly
func (a *ArdroneAdaptor) Finalize() bool {
	return true
}
