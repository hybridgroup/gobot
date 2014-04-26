package gobot

type jsonRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*jsonConnection `json:"connections"`
	Devices     []*jsonDevice     `json:"devices"`
}

type jsonDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *jsonConnection `json:"connection"`
	Commands   []string        `json:"commands"`
}

type jsonConnection struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	Adaptor string `json:"adaptor"`
}
