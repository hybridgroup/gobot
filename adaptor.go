package gobot

type Adaptor struct {
	Name      string
	Port      string
	Connected bool
	Params    map[string]interface{}
}

type AdaptorInterface interface {
	Finalize() bool
	Connect() bool
	port() string
	name() string
	setName(string)
	params() map[string]interface{}
}

func (a *Adaptor) port() string {
	return a.Port
}

func (a *Adaptor) name() string {
	return a.Name
}

func (a *Adaptor) setName(s string) {
	a.Name = s
}

func (a *Adaptor) params() map[string]interface{} {
	return a.Params
}
