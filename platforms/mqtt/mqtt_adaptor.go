package mqtt

import (
	"fmt"
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/hybridgroup/gobot"
	"net/url"
)

type MqttAdaptor struct {
	gobot.Adaptor
	Host   string
	Client *mqtt.MqttClient
}

// NewMqttAdaptor creates a new mqtt adaptor with specified name
func NewMqttAdaptor(name string, host string) *MqttAdaptor {
	return &MqttAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"MqttAdaptor",
		),
		Host: host,
	}
}

// Connect returns true if connection to mqtt is established succesfully
func (a *MqttAdaptor) Connect() bool {
	opts := createClientOptions("sub", a.Host)
	a.Client = mqtt.NewClient(opts)
	a.Client.Start()
	return true
}

// Reconnect retries connection to mqtt. Returns true if succesfull
func (a *MqttAdaptor) Reconnect() bool {
	return true
}

// Disconnect returns true if connection to mqtt is closed succesfully
func (a *MqttAdaptor) Disconnect() bool {
	if a.Client != nil {
		a.Client.Disconnect(500)
	}
	return true
}

// Finalize returns true if connection to mqtt is finalized succesfully
func (a *MqttAdaptor) Finalize() bool {
	a.Disconnect()
	return true
}

func (a *MqttAdaptor) Publish(topic string, message []byte) int {
	a.Client.Publish(0, topic, message)
	return 0
}

func (a *MqttAdaptor) Subscribe(topic string) int {
	return 0
}

func createClientOptions(clientId, raw string) *mqtt.ClientOptions {
	uri, _ := url.Parse(raw)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetClientId(clientId)

	return opts
}
