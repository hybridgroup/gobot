package gobot

import (
	"log"
	"reflect"
)

type connection struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"adaptor"`
	Adaptor AdaptorInterface       `json:"-"`
	Port    string                 `json:"-"`
	Robot   *Robot                 `json:"-"`
	Params  map[string]interface{} `json:"-"`
}

type Connection interface {
	Connect() bool
	Finalize() bool
}

func NewConnection(adaptor AdaptorInterface, r *Robot) *connection {
	c := new(connection)
	c.Type = reflect.ValueOf(adaptor).Type().String()
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
	log.Println("Connecting to " + c.Name + " on port " + c.Port + "...")
	return c.Adaptor.Connect()
}

func (c *connection) Finalize() bool {
	log.Println("Finalizing " + c.Name + "...")
	return c.Adaptor.Finalize()
}
