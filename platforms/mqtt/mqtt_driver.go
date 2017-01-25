package mqtt

import "gobot.io/x/gobot"

// Driver for mqtt
type Driver struct {
	name       string
	topic      string
	connection gobot.Connection
	gobot.Eventer
	gobot.Commander
}

// NewDriver returns a new Gobot MQTT Driver
func NewDriver(a *Adaptor, topic string) *Driver {
	m := &Driver{
		name:       "MQTT",
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
func (m *Driver) Publish(message []byte) bool {
	return m.adaptor().Publish(m.topic, message)
}

// OnData subscribes to the current device topic, and then calls the message handler function when data is received
func (m *Driver) OnData(f func(s []byte)) bool {
	return m.adaptor().On(m.topic, f)
}
