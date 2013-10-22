package gobot

import "fmt"

type Connection struct {
  Connection_id string
  Name string
  Adaptor string
  Port string
  Parent string
}

//func (c *Connection) New(c Connection) *Connection {
//  return c
//}

func (c *Connection) Connect() {
  fmt.Println("Connecting to "+ c.Name + " on port " + c.Port + "...")
}

func (c *Connection) Disconnect() {
  fmt.Println("Diconnecting from "+ c.Name + " on port " + c.Port + "...")
}
