package gobot

type Adaptor struct {
	Name      string                 `json:"name"`
	Port      string                 `json:"port"`
	Connected bool                   `json:"Connected"`
	Params    map[string]interface{} `json:"params"`
}

type AdaptorInterface interface {
	Finalize() bool
	Connect() bool
}
