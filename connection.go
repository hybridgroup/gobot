package gobot

import (
  "fmt"
  "reflect"
)

type Connection struct {
  Name string
  Adaptor interface{}
  Port Port
  Robot *Robot
  Params map[string]string

}

func NewConnection(a interface{}, r *Robot) *Connection {
  c := new(Connection)
  c.Name = reflect.ValueOf(a).Elem().FieldByName("Name").String()
  c.Robot = r
  c.Adaptor = a
  return c
}

func (c *Connection) Connect() {
  fmt.Println("Connecting to " + reflect.ValueOf(c).Elem().FieldByName("Name").String() + " on port " + c.Port.ToString() + "...")
  reflect.ValueOf(c.Adaptor).MethodByName("Connect").Call([]reflect.Value{})
}

func (c *Connection) Disconnect() {
  reflect.ValueOf(c.Adaptor).MethodByName("Disconnect").Call([]reflect.Value{})
}

func (c *Connection) IsConnected() bool {
  return reflect.ValueOf(c.Adaptor).MethodByName("IsConnect").Call([]reflect.Value{})[0].Bool()
}

func (c *Connection) AdaptorName() string {
  return c.Name
}
