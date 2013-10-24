package gobot

import (
  "time"
  "fmt"
  "math/rand"
)

var connectionTypes []*ConnectionType
var deviceTypes []*DeviceType

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
  r.initConnections(r.Connections)
  r.initDevices(r.Devices)
  r.startConnections()
  r.startDevices()
  r.Work()
  for{time.Sleep(1 * time.Second)}
}

func (r *Robot) initConnections(connections []Connection) {
  connectionTypes := make([]*ConnectionType, len(connections))
  fmt.Println("Initializing connections...")
  for i := range connections {
    fmt.Println("Initializing connection " + connections[i].Name + "...")
    connections[i].Robot = *r
    connectionTypes[i] = NewConnection(connections[i])
  }
}

func (r *Robot) initDevices(devices []Device) {
  deviceTypes := make([]*DeviceType, len(devices))
  fmt.Println("Initializing devices...")
  for i := range devices {
    fmt.Println("Initializing device " + devices[i].Name + "...")
    devices[i].Robot = *r
    deviceTypes[i] = NewDevice(devices[i])
  }
}

func (r *Robot) startConnections() {
  fmt.Println("Starting connections...")
  for i := range connectionTypes {
    fmt.Println("Starting connection " + connectionTypes[i].Name + "...")
    connectionTypes[i].Connect()
  }
}

func (r *Robot) startDevices() {
  fmt.Println("Starting devices...")
  for i := range deviceTypes {
    fmt.Println("Starting device " + deviceTypes[i].Name + "...")
    deviceTypes[i].Start()
  }
}
//    # Terminate all connections
//    def disconnect
//      connections.each {|k, c| c.async.disconnect}
//    end

//    # @return [Connection] default connection
//    def default_connection
//      connections.values.first
//    end

//    # @return [Collection] connection types
//    def connection_types
//      current_class.connection_types ||= [{:name => :passthru}]
//    end

//    # @return [Collection] device types
//    def device_types
//      current_class.device_types ||= []
//      current_class.device_types
//    end

//    # @return [Proc] current working code
//    def working_code
//      current_class.working_code ||= proc {puts "No work defined."}
//    end

//    # @param [Symbol] period
//    # @param [Numeric] interval
//    # @return [Boolean] True if there is recurring work for the period and interval
//    def has_work?(period, interval)
//      current_instance.timers.find {|t| t.recurring == (period == :every) && t.interval == interval}
//    end
