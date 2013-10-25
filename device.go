package gobot

import (
  "fmt"
  "strconv"
  "reflect"
)

type Device struct {
  Name string
  Robot *Robot
  Connection *Connection 
  Driver *Driver
  Params map[string]string
}

func NewDevice(d reflect.Value, r *Robot) *Device {
  dt := new(Device)
  dt.Name = reflect.ValueOf(d).FieldByName("Name").String()
  dt.Robot = r
  dt.Driver = new(Driver)
  dt.Driver.Pin = reflect.ValueOf(d).FieldByName("Pin").String()
  dt.Driver.Interval, _ = strconv.ParseFloat(reflect.ValueOf(d).FieldByName("Interval").String(), 64)
  dt.Connection = new(Connection)
  return dt
}

func (dt *Device) Start() {
  fmt.Println("Device " + dt.Name + "started")
  dt.Driver.Start()
}
func (dt *Device) Command(method_name string, arguments []string) {
  //dt.Driver.Command(method_name, arguments)
}
