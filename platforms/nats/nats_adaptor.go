package nats

import (
	"github.com/nats-io/nats"
	"gobot.io/x/gobot"
	"net/url"
	"strings"
)

// Adaptor is a configuration struct for interacting with a NATS server.
// Name is a logical name for the adaptor/nats server connection.
// Host is in the form "localhost:4222" which is the hostname/ip and port of the nats server.
// ClientID is a unique identifier integer that specifies the identity of the client.
type Adaptor struct {
	name     string
	Host     string
	clientID int
	username string
	password string
	client   *nats.Conn
	connect  func() (*nats.Conn, error)
}

// Message is a message received from the server.
type Message *nats.Msg

// NewAdaptor populates a new NATS Adaptor.
func NewAdaptor(host string, clientID int, options ...nats.Option) *Adaptor {
	hosts, err := processHostString(host)

	return &Adaptor{
		name:     gobot.DefaultName("NATS"),
		Host:     hosts,
		clientID: clientID,
		connect: func() (*nats.Conn, error) {
			if err != nil {
				return nil, err
			}
			return nats.Connect(hosts, options...)
		},
	}
}

// NewAdaptorWithAuth populates a NATS Adaptor including username and password.
func NewAdaptorWithAuth(host string, clientID int, username string, password string, options ...nats.Option) *Adaptor {
	hosts, err := processHostString(host)

	return &Adaptor{
		Host:     hosts,
		clientID: clientID,
		username: username,
		password: password,
		connect: func() (*nats.Conn, error) {
			if err != nil {
				return nil, err
			}
			return nats.Connect(hosts, append(options, nats.UserInfo(username, password))...)
		},
	}
}

func processHostString(host string) (string, error) {
	urls := strings.Split(host, ",")
	for i, s := range urls {
		s = strings.TrimSpace(s)
		if !strings.HasPrefix(s, "tls://") && !strings.HasPrefix(s, "nats://") {
			s = "nats://" + s
		}

		u, err := url.Parse(s)
		if err != nil {
			return "", err
		}

		urls[i] = u.String()
	}

	return strings.Join(urls, ","), nil
}

// Name returns the logical client name.
func (a *Adaptor) Name() string { return a.name }

// SetName sets the logical client name.
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect makes a connection to the Nats server.
func (a *Adaptor) Connect() (err error) {
	a.client, err = a.connect()
	return
}

// Disconnect from the nats server.
func (a *Adaptor) Disconnect() (err error) {
	if a.client != nil {
		a.client.Close()
	}
	return
}

// Finalize is simply a helper method for the disconnect.
func (a *Adaptor) Finalize() (err error) {
	a.Disconnect()
	return
}

// Publish sends a message with the particular topic to the nats server.
func (a *Adaptor) Publish(topic string, message []byte) bool {
	if a.client == nil {
		return false
	}
	a.client.Publish(topic, message)
	return true
}

// On is an event-handler style subscriber to a particular topic (named event).
// Supply a handler function to use the bytes returned by the server.
func (a *Adaptor) On(event string, f func(msg Message)) bool {
	if a.client == nil {
		return false
	}
	a.client.Subscribe(event, func(msg *nats.Msg) {
		f(msg)
	})

	return true
}
