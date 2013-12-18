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
	Disconnect() bool
	Reconnect() bool
}
