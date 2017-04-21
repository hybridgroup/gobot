package nats

import "gobot.io/x/gobot"

const (
	// Data event when data is available for Driver
	Data = "data"

	// Error event when error occurs in Driver
	Error = "error"
)

// Driver for NATS
type Driver struct {
	name       string
	topic      string
	connection gobot.Connection
	gobot.Eventer
	gobot.Commander
}

// NewDriver returns a new Gobot NATS Driver
func NewDriver(a *Adaptor, topic string) *Driver {
	m := &Driver{
		name:       gobot.DefaultName("NATS"),
		topic:      topic,
		connection: a,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}

	return m
}

// Name returns name for the Driver
func (m *Driver) Name() string { return m.name }

// Name sets name for the Driver
func (m *Driver) SetName(name string) { m.name = name }

// Connection returns Connections used by the Driver
func (m *Driver) Connection() gobot.Connection {
	return m.connection
}

func (m *Driver) adaptor() *Adaptor {
	return m.Connection().(*Adaptor)
}

// Start starts the Driver
func (m *Driver) Start() error {
	return nil
}

// Halt halts the Driver
func (m *Driver) Halt() error {
	return nil
}

// Topic returns the current topic for the Driver
func (m *Driver) Topic() string { return m.topic }

// SetTopic sets the current topic for the Driver
func (m *Driver) SetTopic(topic string) { m.topic = topic }

// Publish a message to the current device topic
func (m *Driver) Publish(data interface{}) bool {
	message := data.([]byte)
	return m.adaptor().Publish(m.topic, message)
}

// On subscribes to data updates for the current device topic,
// and then calls the message handler function when data is received
func (m *Driver) On(n string, f func(msg Message)) error {
	// TODO: also be able to subscribe to Error updates
	m.adaptor().On(m.topic, f)
	return nil
}
