package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"gobot.io/x/gobot"

	paho "github.com/eclipse/paho.mqtt.golang"
	multierror "github.com/hashicorp/go-multierror"
)

// Message is a message received from the broker.
type Message paho.Message

// Adaptor is the Gobot Adaptor for MQTT
type Adaptor struct {
	name          string
	Host          string
	clientID      string
	username      string
	password      string
	useSSL        bool
	serverCert    string
	clientCert    string
	clientKey     string
	autoReconnect bool
	cleanSession  bool
	client        paho.Client
}

// NewAdaptor creates a new mqtt adaptor with specified host and client id
func NewAdaptor(host string, clientID string) *Adaptor {
	return &Adaptor{
		name:          gobot.DefaultName("MQTT"),
		Host:          host,
		autoReconnect: false,
		cleanSession:  true,
		useSSL:        false,
		clientID:      clientID,
	}
}

// NewAdaptorWithAuth creates a new mqtt adaptor with specified host, client id, username, and password.
func NewAdaptorWithAuth(host, clientID, username, password string) *Adaptor {
	return &Adaptor{
		name:          "MQTT",
		Host:          host,
		autoReconnect: false,
		cleanSession:  true,
		useSSL:        false,
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

// CleanSession returns the MQTT CleanSession setting
func (a *Adaptor) CleanSession() bool { return a.cleanSession }

// SetCleanSession sets the MQTT CleanSession setting. Should be false if reconnect is enabled. Otherwise all subscriptions will be lost
func (a *Adaptor) SetCleanSession(val bool) { a.cleanSession = val }

// UseSSL returns the MQTT server SSL preference
func (a *Adaptor) UseSSL() bool { return a.useSSL }

// SetUseSSL sets the MQTT server SSL preference
func (a *Adaptor) SetUseSSL(val bool) { a.useSSL = val }

// ServerCert returns the MQTT server SSL cert file
func (a *Adaptor) ServerCert() string { return a.serverCert }

// SetServerCert sets the MQTT server SSL cert file
func (a *Adaptor) SetServerCert(val string) { a.serverCert = val }

// ClientCert returns the MQTT client SSL cert file
func (a *Adaptor) ClientCert() string { return a.clientCert }

// SetClientCert sets the MQTT server SSL cert file
func (a *Adaptor) SetClientCert(val string) { a.clientCert = val }

// ClientKey returns the MQTT client SSL key file
func (a *Adaptor) ClientKey() string { return a.clientKey }

// SetClientKey sets the MQTT client SSL key file
func (a *Adaptor) SetClientKey(val string) { a.clientKey = val }

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
func (a *Adaptor) On(event string, f func(msg Message)) bool {
	if a.client == nil {
		return false
	}
	a.client.Subscribe(event, 0, func(client paho.Client, msg paho.Message) {
		f(msg)
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
	opts.CleanSession = a.cleanSession

	if a.UseSSL() {
		opts.SetTLSConfig(a.newTLSConfig())
	}
	return opts
}

// newTLSConfig sets the TLS config in the case that we are using
// an MQTT broker with TLS
func (a *Adaptor) newTLSConfig() *tls.Config {
	// Import server certificate
	var certpool *x509.CertPool
	if len(a.ServerCert()) > 0 {
		certpool = x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(a.ServerCert())
		if err == nil {
			certpool.AppendCertsFromPEM(pemCerts)
		}
	}

	// Import client certificate/key pair
	var certs []tls.Certificate
	if len(a.ClientCert()) > 0 && len(a.ClientKey()) > 0 {
		cert, err := tls.LoadX509KeyPair(a.ClientCert(), a.ClientKey())
		if err != nil {
			// TODO: proper error handling
			panic(err)
		}
		certs = append(certs, cert)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: certs,
	}
}
