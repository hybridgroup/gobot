package gpio

import (
	"math"
	"time"

	"gobot.io/x/gobot"
)

var _ gobot.Driver = (*GroveDHT11SensorDriver)(nil)

// GroveDHT11SensorOption options pattern for the GroveDHT11Sensor constructor function.
type GroveDHT11SensorOption func(g *GroveDHT11SensorDriver)

// NewGroveDHT11SensorDriver returns a new NewGroveDHT11SensorDriver with default values given an DHTReader and pin.
//
// Optionally acceptes zero or more GroveDHT11SensorOption.
// According specification (see https://www.mouser.com/datasheet/2/758/DHT11-Technical-Data-Sheet-Translated-Version-1143054.pdf)
// sampling interval must be greater than or equal 1 second. Default interval is one second.
func NewGroveDHT11SensorDriver(d DHTReader, pin string, opts ...GroveDHT11SensorOption) *GroveDHT11SensorDriver {
	g := &GroveDHT11SensorDriver{
		name:       gobot.DefaultName("GroveDHT11Sensor"),
		connection: d,
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		halt:       make(chan bool),
		interval:   1000 * time.Millisecond,
		state:      make(chan GroveDHT11SensorState, 1),
	}

	g.state <- GroveDHT11SensorState{Temperature: 0, Humidity: 0}

	for _, fn := range opts {
		fn(g)
	}

	g.AddEvent(Data)
	g.AddEvent(Error)

	return g
}

// WithGroveDHT11SensorInterval configures the driver to use given duration as polling interval.
//
// According specification (see https://www.mouser.com/datasheet/2/758/DHT11-Technical-Data-Sheet-Translated-Version-1143054.pdf)
// sampling interval must be greater than or equal 1 second. Default interval is one second.
func WithGroveDHT11SensorInterval(duration time.Duration) GroveDHT11SensorOption {
	return func(g *GroveDHT11SensorDriver) {
		if duration >= time.Millisecond*1000 {
			g.interval = duration
		}
	}
}

// GroveDHT11SensorDriver driver for the Grove temperature and humidity sensor DHT11.
type GroveDHT11SensorDriver struct {
	name       string
	pin        string
	halt       chan bool
	connection DHTReader
	interval   time.Duration
	state      chan GroveDHT11SensorState
	gobot.Eventer
}

// Temperature returns last read temperature in Celsius.
// This function is safe to use in multiple goroutines.
func (g GroveDHT11SensorDriver) Temperature() (val float32) {
	state := <-g.state
	defer func() { g.state <- state }()
	return state.Temperature
}

// Humidity returns last read humidity.
// This function is safe to use in multiple goroutines.
func (g GroveDHT11SensorDriver) Humidity() (val float32) {
	state := <-g.state
	defer func() { g.state <- state }()
	return state.Humidity
}

// Name returns the label for the Driver
func (g GroveDHT11SensorDriver) Name() string {
	return g.name
}

// SetName sets the label for the Driver
func (g *GroveDHT11SensorDriver) SetName(name string) {
	g.name = name
}

// Pin returns the configured pin which this sensor is connected to.
func (g GroveDHT11SensorDriver) Pin() string {
	return g.pin
}

// Start initiates the Driver
func (g *GroveDHT11SensorDriver) Start() error {
	go func() {
		for {
			// Read values
			t, h, err := g.connection.ReadDHT(g.pin)

			if err != nil {
				g.Publish(Error, err)
			} else {
				state := <-g.state
				state.Temperature = float32(math.Round(float64(t)))
				state.Humidity = float32(math.Round(float64(h)))
				g.state <- state
				g.Publish(Data, state)
			}

			select {
			case <-time.After(g.interval):
			case <-g.halt:
				return
			}
		}
	}()

	return nil
}

// Halt terminates the Driver
func (g GroveDHT11SensorDriver) Halt() error {
	g.halt <- true
	return nil
}

// Connection returns the Connection associated with the Driver
func (g GroveDHT11SensorDriver) Connection() gobot.Connection {
	return g.connection.(gobot.Connection)
}

// GroveDHT11SensorState represents current measured temperature and humidity of the DHT11 sensor.
type GroveDHT11SensorState struct {
	Temperature float32
	Humidity    float32
}
