package nats

import (
	"github.com/nats-io/nats"
)

// NatsAdaptor is a configuration struct for interacting with a nats server.
// Name is a logical name for the adaptor/nats server connection.
// Host is in the form "localhost:4222" which is the hostname/ip and port of the nats server.
// ClientID is a unique identifier integer that specifies the identity of the client.
type NatsAdaptor struct {
	name     string
	Host     string
	clientID int
	username string
	password string
	client   *nats.Conn
}

// NewNatsAdaptor populates a new NatsAdaptor.
func NewNatsAdaptor(name string, host string, clientID int) *NatsAdaptor {
	return &NatsAdaptor{
		name:     name,
		Host:     host,
		clientID: clientID,
	}
}

// NewNatsAdaptorWithAuth populates a NatsAdaptor including username and password.
func NewNatsAdaptorWithAuth(name string, host string, clientID int, username string, password string) *NatsAdaptor {
	return &NatsAdaptor{
		name:     name,
		Host:     host,
		clientID: clientID,
		username: username,
		password: password,
	}
}

// Name returns the logical client name.
func (a *NatsAdaptor) Name() string { return a.name }

// Connect makes a connection to the Nats server.
func (a *NatsAdaptor) Connect() (errs []error) {

	auth := ""
	if a.username != "" && a.password != "" {
		auth = a.username + ":" + a.password + "@"
	}

	defaultURL := "nats://" + auth + a.Host

	var err error
	a.client, err = nats.Connect(defaultURL)
	if err != nil {
		return append(errs, err)
	}
	return
}

// Disconnect from the nats server. Returns an error if the client doesn't exist.
func (a *NatsAdaptor) Disconnect() (err error) {
	if a.client != nil {
		a.client.Close()
	}
	return
}

// Finalize is simply a helper method for the disconnect.
func (a *NatsAdaptor) Finalize() (errs []error) {
	a.Disconnect()
	return
}

// Publish sends a message with the particular topic to the nats server.
func (a *NatsAdaptor) Publish(topic string, message []byte) bool {
	if a.client == nil {
		return false
	}
	a.client.Publish(topic, message)
	return true
}

// On is an event-handler style subscriber to a particular topic (named event).
// Supply a handler function to use the bytes returned by the server.
func (a *NatsAdaptor) On(event string, f func(s []byte)) bool {
	if a.client == nil {
		return false
	}
	a.client.Subscribe(event, func(msg *nats.Msg) {
		f(msg.Data)
	})
	return true
}
