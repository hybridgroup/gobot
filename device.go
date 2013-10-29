package gobot

import (
  "fmt"
  "reflect"
)

type Device struct {
  Name string
  Interval string
  Robot *Robot
  Connection *Connection 
  Driver interface{}
  Params map[string]string
}

func NewDevice(driver interface{}, r *Robot) *Device {
  d := new(Device)
  d.Name = reflect.ValueOf(driver).Elem().FieldByName("Name").String()
  d.Robot = r
  d.Driver = driver
  d.Connection = new(Connection)
  return d
}

func (d *Device) Start() {
  fmt.Println("Device " + d.Name + " started")
  reflect.ValueOf(d.Driver).MethodByName("StartDriver").Call([]reflect.Value{})
}

func (d *Device) Command(method_name string, arguments []string) {
  //dt.Driver.Command(method_name, arguments)
}
