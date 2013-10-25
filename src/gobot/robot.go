package gobot

import (
  "time"
  "fmt"
  "math/rand"
  "reflect"
)

var connections []*Connection
var devices []*Device

type Robot struct {
  Connections []interface{}
  Devices []interface{}
  Name string
  Work func()
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
  connections := make([]*Connection, len(r.Connections))
  fmt.Println("Initializing connections...")
  for i := range r.Connections {
    fmt.Println("Initializing connection " + reflect.ValueOf(r.Connections[i]).Elem().FieldByName("Name").String() + "...")
    connections[i] = NewConnection(reflect.ValueOf(r.Connections[i]).Elem().FieldByName("Adaptor"), r)
  }
}

func (r *Robot) initDevices() {
  devices := make([]*Device, len(r.Devices))
  fmt.Println("Initializing devices...")
  for i := range r.Devices {
    fmt.Println("Initializing device " + reflect.ValueOf(r.Devices[i]).Elem().FieldByName("Name").String() + "...")
    devices[i] = NewDevice(reflect.ValueOf(r.Connections[i]).Elem().FieldByName("Driver"), r)
  }
}

func (r *Robot) startConnections() {
  fmt.Println("Starting connections...")
  for i := range connections {
    fmt.Println("Starting connection " + connections[i].Name + "...")
    connections[i].Connect()
  }
}

func (r *Robot) startDevices() {
  fmt.Println("Starting devices...")
  for i := range devices {
    fmt.Println("Starting device " + devices[i].Name + "...")
    devices[i].Start()
  }
}
