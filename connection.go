package gobot

import (
	"fmt"
)

type connection struct {
	Name    string
	Adaptor AdaptorInterface
	Port    string                 `json:"-"`
	Robot   *Robot                 `json:"-"`
	Params  map[string]interface{} `json:"-"`
}

type Connection interface {
	Connect() bool
	Disconnect() bool
	Finalize() bool
	Reconnect() bool
}

func NewConnection(adaptor AdaptorInterface, r *Robot) *connection {
	c := new(connection)
	c.Name = FieldByNamePtr(adaptor, "Name").String()
	c.Port = FieldByNamePtr(adaptor, "Port").String()
	c.Params = make(map[string]interface{})
	keys := FieldByNamePtr(adaptor, "Params").MapKeys()
	for k := range keys {
		c.Params[keys[k].String()] = FieldByNamePtr(adaptor, "Params").MapIndex(keys[k])
	}
	c.Robot = r
	c.Adaptor = adaptor
	return c
}

func (c *connection) Connect() bool {
	fmt.Println("Connecting to " + c.Name + " on port " + c.Port + "...")
	return c.Adaptor.Connect()
}

func (c *connection) Disconnect() bool {
	return c.Adaptor.Disconnect()
}

func (c *connection) Finalize() bool {
	return c.Adaptor.Finalize()
}

func (c *connection) Reconnect() bool {
	return c.Adaptor.Reconnect()
}

func (c *connection) AdaptorName() string {
	return c.Name
}
