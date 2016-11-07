package mqtt

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	multierror "github.com/hashicorp/go-multierror"
)

type Adaptor struct {
	name     string
	Host     string
	clientID string
	username string
	password string
	client   paho.Client
}

// NewAdaptor creates a new mqtt adaptor with specified host and client id
func NewAdaptor(host string, clientID string) *Adaptor {
	return &Adaptor{
		name:     "MQTT",
		Host:     host,
		clientID: clientID,
	}
}

func NewAdaptorWithAuth(host, clientID, username, password string) *Adaptor {
	return &Adaptor{
		name:     "MQTT",
		Host:     host,
		clientID: clientID,
		username: username,
		password: password,
	}
}

func (a *Adaptor) Name() string     { return a.name }
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect returns true if connection to mqtt is established
func (a *Adaptor) Connect() (err error) {
	a.client = paho.NewClient(createClientOptions(a.clientID, a.Host, a.username, a.password))
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

// Subscribe to a topic, and then call the message handler function when data is received
func (a *Adaptor) On(event string, f func(s []byte)) bool {
	if a.client == nil {
		return false
	}
	a.client.Subscribe(event, 0, func(client paho.Client, msg paho.Message) {
		f(msg.Payload())
	})
	return true
}

func createClientOptions(clientId, raw, username, password string) *paho.ClientOptions {
	opts := paho.NewClientOptions()
	opts.AddBroker(raw)
	opts.SetClientID(clientId)
	if username != "" && password != "" {
		opts.SetPassword(password)
		opts.SetUsername(username)
	}
	opts.AutoReconnect = false
	return opts
}
