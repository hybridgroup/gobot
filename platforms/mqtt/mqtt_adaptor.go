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
	client *mqtt.MqttClient
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

// Connect returns true if connection to mqtt is established
func (a *MqttAdaptor) Connect() bool {
	opts := createClientOptions("sub", a.Host)
	a.client = mqtt.NewClient(opts)
	a.client.Start()
	return true
}

// Reconnect retries connection to mqtt. Returns true if successful
func (a *MqttAdaptor) Reconnect() bool {
	return true
}

// Disconnect returns true if connection to mqtt is closed
func (a *MqttAdaptor) Disconnect() bool {
	if a.client != nil {
		a.client.Disconnect(500)
	}
	return true
}

// Finalize returns true if connection to mqtt is finalized succesfully
func (a *MqttAdaptor) Finalize() bool {
	a.Disconnect()
	return true
}

func (a *MqttAdaptor) Publish(topic string, message []byte) int {
	m := mqtt.NewMessage(message)
	a.client.PublishMessage(topic, m)
	return 0
}

func (a *MqttAdaptor) On(event string, f func(s interface{})) {
	t, _ := mqtt.NewTopicFilter(event, 0)
	a.client.StartSubscription(func(client *mqtt.MqttClient, msg mqtt.Message) {
		f(msg.Payload())
	}, t)
}

func createClientOptions(clientId, raw string) *mqtt.ClientOptions {
	uri, _ := url.Parse(raw)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetClientId(clientId)

	return opts
}
