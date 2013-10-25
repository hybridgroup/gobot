package gobot

import (
  "fmt"
  "reflect"
)

type Connection struct {
  Name string
  Adaptor *Adaptor
  Port Port
  Robot *Robot
  Params map[string]string
}

func NewConnection(a reflect.Value, r *Robot) *Connection {
  c := new(Connection)
  c.Name = reflect.ValueOf(a).FieldByName("Name").String()
  c.Robot = r
  c.Adaptor = new(Adaptor)
  c.Adaptor.Name = reflect.ValueOf(a).FieldByName("Name").String()
  return c
}

func (c *Connection) Connect() {
  fmt.Println("Connecting to "+ c.Adaptor.Name + " on port " + c.Port.ToString() + "...")
  c.Adaptor.Connect()
}

func (c *Connection) Disconnect() {
  fmt.Println("Diconnecting from "+ c.Adaptor.Name + " on port " + c.Port.ToString() + "...")
  c.Adaptor.Disconnect()
}

func (c *Connection) IsConnected() bool {
  return c.Adaptor.IsConnected()
}

func (c *Connection) AdaptorName() string {
  return c.Adaptor.Name
}
