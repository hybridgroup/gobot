package ardrone

import (
	client "github.com/hybridgroup/go-ardrone/client"
	"github.com/hybridgroup/gobot"
)

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

func (a *ArdroneAdaptor) Connect() bool {
	a.connect(a)
	return true
}

func (a *ArdroneAdaptor) Finalize() bool {
	return true
}
