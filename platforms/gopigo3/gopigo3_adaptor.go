package gopigo3

import (
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/hybridgroup/gobot/drivers/spi"
	go_spi "golang.org/x/exp/io/spi"
)

const (
	GOPIGO_ADDRESS = 8 // The default GoPiGo address
)

// Adaptor represents a connection to a GoPiGo3
type Adaptor struct {
	name  string
	raspi *raspi.Adaptor
	connect func() *spi.Connection
	connection *spi.Connection
}

// NewAdaptor creates and returns a new GoPiGo adaptor
func NewAdaptor() (*Adaptor, error) {
	a := &Adaptor{
		name:  "GoPiGo",
		raspi: raspi.NewAdaptor(),
		connect: func() (*spi.Connection) {
			return spi.NewConnection(&go_spi.Devfs{
				Dev:      "/dev/spidev0.1",
				Mode:     go_spi.Mode0,
				MaxSpeed: 500000,
			}, GOPIGO_ADDRESS)
		},
	}

	return a, nil
}

// Name returns the Adaptor's name
func (a *Adaptor) Name() string {
	return a.name
}

// SetName sets the Adaptor's name
func (a *Adaptor) SetName(name string) {
	a.name = name
}

// Connect makes a connection to the GoPiGo3
func (a *Adaptor) Connect() error {
	conn := a.connect()
	a.connection = conn

	return nil
}

// Finalize closes the connection to the GoPiGo3
func (a *Adaptor) Finalize() error {
	return a.connection.Close()
}
