package mqtt

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	multierror "github.com/hashicorp/go-multierror"
)

// Adaptor is the Gobot Adaptor for MQTT
type Adaptor struct {
	name          string
	Host          string
	clientID      string
	username      string
	password      string
	autoReconnect bool
	client        paho.Client
}

// NewAdaptor creates a new mqtt adaptor with specified host and client id
func NewAdaptor(host string, clientID string) *Adaptor {
	return &Adaptor{
		name:          "MQTT",
		Host:          host,
		autoReconnect: false,
		clientID:      clientID,
	}
}

// NewAdaptorWithAuth creates a new mqtt adaptor with specified host, client id, username, and password.
func NewAdaptorWithAuth(host, clientID, username, password string) *Adaptor {
	return &Adaptor{
		name:          "MQTT",
		Host:          host,
		autoReconnect: false,
		clientID:      clientID,
		username:      username,
		password:      password,
	}
}

// Name returns the MQTT Adaptor's name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the MQTT Adaptor's name
func (a *Adaptor) SetName(n string) { a.name = n }

// Port returns the Host name
func (a *Adaptor) Port() string { return a.Host }

// AutoReconnect returns the MQTT AutoReconnect setting
func (a *Adaptor) AutoReconnect() bool { return a.autoReconnect }

// SetAutoReconnect sets the MQTT AutoReconnect setting
func (a *Adaptor) SetAutoReconnect(val bool) { a.autoReconnect = val }

// Connect returns true if connection to mqtt is established
func (a *Adaptor) Connect() (err error) {
	a.client = paho.NewClient(a.createClientOptions())
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		err = multierror.Append(err, token.Error())
	}

	return
}

// Disconnect returns true if connection to mqtt is closed
func (a *Adaptor) Disconnect() (err error) {
	if a.client != nil {
		a.client.Disconnect(500)
	}
	return
}

// Finalize returns true if connection to mqtt is finalized successfully
func (a *Adaptor) Finalize() (err error) {
	a.Disconnect()
	return
}

// Publish a message under a specific topic
func (a *Adaptor) Publish(topic string, message []byte) bool {
	if a.client == nil {
		return false
	}
	a.client.Publish(topic, 0, false, message)
	return true
}

// On subscribes to a topic, and then calls the message handler function when data is received
func (a *Adaptor) On(event string, f func(s []byte)) bool {
	if a.client == nil {
		return false
	}
	a.client.Subscribe(event, 0, func(client paho.Client, msg paho.Message) {
		f(msg.Payload())
	})
	return true
}

func (a *Adaptor) createClientOptions() *paho.ClientOptions {
	opts := paho.NewClientOptions()
	opts.AddBroker(a.Host)
	opts.SetClientID(a.clientID)
	if a.username != "" && a.password != "" {
		opts.SetPassword(a.password)
		opts.SetUsername(a.username)
	}
	opts.AutoReconnect = a.autoReconnect
	return opts
}
