package gobot

import "fmt"

type Adaptor struct {
	name        string
	port        string
	connected   bool
	params      map[string]interface{}
	adaptorType string
}

type AdaptorInterface interface {
	Finalize() bool
	Connect() bool
	Port() string
	Name() string
	Type() string
	Connected() bool
	SetConnected(bool)
	SetName(string)
	Params() map[string]interface{}
	ToJSON() *JSONConnection
}

func NewAdaptor(name string, adaptorType string, v ...interface{}) *Adaptor {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	a := &Adaptor{
		adaptorType: adaptorType,
		name:        name,
		port:        "",
		params:      make(map[string]interface{}),
	}

	for i := range v {
		switch v[i].(type) {
		case string:
			a.port = v[i].(string)
		case map[string]interface{}:
			a.params = v[i].(map[string]interface{})
		default:
			fmt.Println("Unknown argument passed to NewAdaptor")
		}
	}

	return a
}

func (a *Adaptor) Port() string {
	return a.port
}

func (a *Adaptor) Name() string {
	return a.name
}

func (a *Adaptor) SetName(s string) {
	a.name = s
}

func (a *Adaptor) Type() string {
	return a.adaptorType
}

func (a *Adaptor) Connected() bool {
	return a.connected
}

func (a *Adaptor) SetConnected(b bool) {
	a.connected = b
}

func (a *Adaptor) Params() map[string]interface{} {
	return a.params
}

func (a *Adaptor) ToJSON() *JSONConnection {
	return &JSONConnection{
		Name:    a.Name(),
		Port:    a.Port(),
		Adaptor: a.Type(),
	}
}
