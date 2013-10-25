package gobot

import "fmt"

type Driver struct {
  Interval float64
  Pin string
  Name string
  //Robot Robot
  Params map[string]string
}

func NewDriver(d Driver) Driver {
  return d
}

// @return [Connection] parent connection
func (d *Driver) Connection() *interface{}{
  //return d.Robot.Connections[0]
  return new(interface{})
}

// @return [String] parent pin
//func (d *Driver) Pin() string {
//  return d.Robot.Devices[0].Pin
//}

// @return [String] parent interval
//func (d *Driver) Interval() string {
//  return d.Robot.Devices[0].Interval
//}

// Generic driver start
func (d *Driver) Start() {
  fmt.Println("Starting driver " +  d.Name + "...")
}

// @return [String] parent topic name
//func eventTopicName(event) {
//  parent.event_topic_name(event)
//}

// @return [Collection] commands
//func commands() {
//  self.class.const_get('COMMANDS')
//}

// Execute command
// @param [Symbol] method_name
// @param [Array]  arguments
//func command(method_name, *arguments) {
//  known_command?(method_name)
//  if arguments.first
//    self.send(method_name, *arguments)
//  else
//    self.send(method_name)
//  end
//  rescue Exception => e
//    Logger.error e.message
//    Logger.error e.backtrace.inspect
//  return nil
//}

// @return [Boolean] True if command exists
//func isKnownCommand(method_name) {
//  return true if commands.include?(method_name.intern)
//
//  Logger.warn("Calling unknown command '#{method_name}'...")
//  return false
//}
