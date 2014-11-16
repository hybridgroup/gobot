package mqtt

import (
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*MqttAdaptor)(nil)

type MqttAdaptor struct {
	gobot.Adaptor
	Host     string
	clientID string
	client   *mqtt.MqttClient
}

// NewMqttAdaptor creates a new mqtt adaptor with specified name, host and client id
func NewMqttAdaptor(name string, host string, clientID string) *MqttAdaptor {
	return &MqttAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"MqttAdaptor",
		),
		Host:     host,
		clientID: clientID,
	}
}

// Connect returns true if connection to mqtt is established
func (a *MqttAdaptor) Connect() error {
	opts := createClientOptions(a.clientID, a.Host)
	a.client = mqtt.NewClient(opts)
	a.client.Start()
	return nil
}

// Disconnect returns true if connection to mqtt is closed
func (a *MqttAdaptor) Disconnect() error {
	if a.client != nil {
		a.client.Disconnect(500)
	}
	return nil
}

// Finalize returns true if connection to mqtt is finalized succesfully
func (a *MqttAdaptor) Finalize() error {
	a.Disconnect()
	return nil
}

// Publish a message under a specific topic
func (a *MqttAdaptor) Publish(topic string, message []byte) bool {
	if a.client == nil {
		return false
	}
	m := mqtt.NewMessage(message)
	a.client.PublishMessage(topic, m)
	return true
}

// Subscribe to a topic, and then call the message handler function when data is received
func (a *MqttAdaptor) On(event string, f func(s []byte)) bool {
	if a.client == nil {
		return false
	}
	t, _ := mqtt.NewTopicFilter(event, 0)
	a.client.StartSubscription(func(client *mqtt.MqttClient, msg mqtt.Message) {
		f(msg.Payload())
	}, t)
	return true
}

func createClientOptions(clientId, raw string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(raw)
	opts.SetClientId(clientId)

	return opts
}
