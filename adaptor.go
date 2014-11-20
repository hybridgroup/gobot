package gobot

import "fmt"

type Adaptor struct {
	name        string
	port        string
	connected   bool
	adaptorType string
}

// AdaptorInterface defines behaviour expected for a Gobot Adaptor
type AdaptorInterface interface {
	Finalize() []error
	Connect() []error
	Port() string
	Name() string
	Type() string
	Connected() bool
	SetConnected(bool)
	SetName(string)
	SetPort(string)
	ToJSON() *JSONConnection
}

// NewAdaptor returns a new Gobot Adaptor
func NewAdaptor(name string, adaptorType string, v ...interface{}) *Adaptor {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	a := &Adaptor{
		adaptorType: adaptorType,
		name:        name,
		port:        "",
	}

	for i := range v {
		switch v[i].(type) {
		case string:
			a.port = v[i].(string)
		}
	}

	return a
}

// Port returns adaptor port
func (a *Adaptor) Port() string {
	return a.port
}

// SetPort sets adaptor port
func (a *Adaptor) SetPort(s string) {
	a.port = s
}

// Name returns adaptor name
func (a *Adaptor) Name() string {
	return a.name
}

// SetName sets adaptor name
func (a *Adaptor) SetName(s string) {
	a.name = s
}

// Type returns adaptor type
func (a *Adaptor) Type() string {
	return a.adaptorType
}

// Connected returns true if adaptor is connected
func (a *Adaptor) Connected() bool {
	return a.connected
}

// SetConnected sets adaptor as connected/disconnected
func (a *Adaptor) SetConnected(b bool) {
	a.connected = b
}

// ToJSON returns a json representation of adaptor
func (a *Adaptor) ToJSON() *JSONConnection {
	return &JSONConnection{
		Name:    a.Name(),
		Adaptor: a.Type(),
	}
}
