package mqtt

import (
  "github.com/hybridgroup/gobot"
)

type MqttAdaptor struct {
  gobot.Adaptor
}

// NewMqttAdaptor creates a new mqtt adaptor with specified name
func NewMqttAdaptor(name string) *MqttAdaptor {
  return &MqttAdaptor{
    Adaptor: *gobot.NewAdaptor(
      name,
      "MqttAdaptor",
    ),
  }
}

// Connect returns true if connection to mqtt is established succesfully
func (a *MqttAdaptor) Connect() bool {
  return true
}

// Reconnect retries connection to mqtt. Returns true if succesfull
func (a *MqttAdaptor) Reconnect() bool {
  return true
}

// Disconnect returns true if connection to mqtt is closed succesfully
func (a *MqttAdaptor) Disconnect() bool {
  return true
}

// Finalize returns true if connection to mqtt is finalized succesfully
func (a *MqttAdaptor) Finalize() bool {
  return true
}

func (s *MqttAdaptor) Publish(topic string, message []byte) int {
  return 0
}

func (s *MqttAdaptor) Subscribe(topic string) int {
  return 0
}
