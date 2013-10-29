package gobot

import (
  "time"
  "fmt"
  "math/rand"
  "reflect"
)

type Robot struct {
  Connections []interface{}
  Devices []interface{}
  Name string
  Work func()
  connections []*Connection
  devices []*Device
}

func (r *Robot) Start() {
  if r.Name == "" {
    rand.Seed( time.Now().UTC().UnixNano())
    i := rand.Int()
    r.Name = fmt.Sprintf("Robot %v", i)
  }
  r.initConnections()
  r.initDevices()
  r.startConnections()
  r.startDevices()
  r.Work()
  for{time.Sleep(10 * time.Millisecond)}
}

func (r *Robot) initConnections() {
  r.connections = make([]*Connection, len(r.Connections))
  fmt.Println("Initializing connections...")
  for i := range r.Connections {
    fmt.Println("Initializing connection " + reflect.ValueOf(r.Connections[i]).Elem().FieldByName("Name").String() + "...")
    r.connections[i] = NewConnection(r.Connections[i], r)
  }
}

func (r *Robot) initDevices() {
  r.devices = make([]*Device, len(r.Devices))
  fmt.Println("Initializing devices...")
  for i := range r.Devices {
    fmt.Println("Initializing device " + reflect.ValueOf(r.Devices[i]).Elem().FieldByName("Name").String() + "...")
    r.devices[i] = NewDevice(r.Devices[i], r)
  }
}

func (r *Robot) startConnections() {
  fmt.Println("Starting connections...")
  for i := range r.connections {
    fmt.Println("Starting connection " + r.connections[i].Name + "...")
    r.connections[i].Connect()
  }
}

func (r *Robot) startDevices() {
  fmt.Println("Starting devices...")
  for i := range r.devices {
    fmt.Println("Starting device " + r.devices[i].Name + "...")
    r.devices[i].Start()
  }
}
