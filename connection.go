package gobot

import (
	"fmt"
	"reflect"
)

type Connection struct {
	Name    string
	Adaptor interface{}
	Port    string
	Robot   *Robot `json:"-"`
	Params  map[string]interface{}
}

func NewConnection(a interface{}, r *Robot) *Connection {
	c := new(Connection)
	c.Name = reflect.ValueOf(a).Elem().FieldByName("Name").String()
	c.Port = reflect.ValueOf(a).Elem().FieldByName("Port").String()
	c.Params = make(map[string]interface{})
	keys := reflect.ValueOf(a).Elem().FieldByName("Params").MapKeys()
	for k := range keys {
		c.Params[keys[k].String()] = reflect.ValueOf(a).Elem().FieldByName("Params").MapIndex(keys[k])
	}
	c.Robot = r
	c.Adaptor = a
	return c
}

func (c *Connection) Connect() {
	fmt.Println("Connecting to " + c.Name + " on port " + c.Port + "...")
	reflect.ValueOf(c.Adaptor).MethodByName("Connect").Call([]reflect.Value{})
}

func (c *Connection) Disconnect() {
	reflect.ValueOf(c.Adaptor).MethodByName("Disconnect").Call([]reflect.Value{})
}

func (c *Connection) IsConnected() bool {
	return reflect.ValueOf(c.Adaptor).MethodByName("IsConnected").Call([]reflect.Value{})[0].Bool()
}

func (c *Connection) AdaptorName() string {
	return c.Name
}
