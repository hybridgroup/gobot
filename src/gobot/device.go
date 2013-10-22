package gobot

import "fmt"

type Device struct {
  Name string
  Pin string
  Parent string
  Connection string
  Interval string
  Driver string
}

//func (d *Device) New() *Device{
//  return d
//}

func (d *Device) Start() {
  fmt.Println("Device " + d.Name + "started")
}
    
func (d *Device) determineConnection(c Connection){
  //d.Parent.connections(c) if c
}

func (d *Device) defaultConnection() {
  //d.Parent.connections.first
}

func requireDriver(driverName string) {
  fmt.Println("dynamic load driver" + driverName)
}
