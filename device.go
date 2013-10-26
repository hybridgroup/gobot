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
  Driver *Driver
  Params map[string]string
}

func NewDevice(d interface{}, r *Robot) *Device {
  dt := new(Device)
  dt.Name = reflect.ValueOf(d).Elem().FieldByName("Name").String()
  dt.Robot = r
  dt.Driver = new(Driver)
  dt.Driver.Pin = reflect.ValueOf(d).Elem().FieldByName("Pin").String()
  dt.Driver.Interval = reflect.ValueOf(d).Elem().FieldByName("Interval").String()
  dt.Driver.Name = reflect.ValueOf(d).Elem().FieldByName("Name").String()
  dt.Connection = new(Connection)
  return dt
}

func (dt *Device) Start() {
  fmt.Println("Device " + dt.Name + " started")
  dt.Driver.Start()
}
func (dt *Device) Command(method_name string, arguments []string) {
  //dt.Driver.Command(method_name, arguments)
}
