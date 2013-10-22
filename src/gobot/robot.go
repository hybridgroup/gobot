package gobot

import (
  "time"
  "fmt"
  "math/rand"
)

var connectionTypes []Connection
var deviceTypes []Device

type Robot struct {
  Connections []Connection
  Devices []Device
  Name string
  Work func()
}

func (r *Robot) Start() {
  if r.Name == "" {
    rand.Seed( time.Now().UTC().UnixNano())
    i := rand.Int()
    r.Name = fmt.Sprintf("Robot %v", i)
  }
  initConnections(r.Connections)
  initDevices(r.Devices)
  startConnections()
  startDevices()
  r.Work()
  for{time.Sleep(1 * time.Second)}
}

func initConnections(connections []Connection) {
  connectionTypes := make([]Connection, len(connections))
  fmt.Println("Initializing connections...")
  for i := range connections {
    fmt.Println("Initializing connection " + connections[i].Name + "...")
//    connectionTypes[i] = Connection.New(connections[i])
    connectionTypes[i] = connections[i]
  }
}

func initDevices(devices []Device) {
  deviceTypes := make([]Device, len(devices))
  fmt.Println("Initializing devices...")
  for i := range devices {
    fmt.Println("Initializing donnection " + devices[i].Name + "...")
//    deviceTypes[i] = Device.New(devices[i])
    deviceTypes[i] = devices[i]
  }
}

func startConnections() {
  fmt.Println("Starting connections...")
  for i := range connectionTypes {
    fmt.Println("Starting connection " + connectionTypes[i].Name + "...")
    connectionTypes[i].Connect()
  }
}

func startDevices() {
  fmt.Println("Starting devices...")
  for i := range deviceTypes {
    fmt.Println("Starting devices " + deviceTypes[i].Name + "...")
    deviceTypes[i].Start()
  }
}
