package gobot

import (
	"errors"
	"fmt"
	"log"
)

// JSONRobot a JSON representation of a robot.
type JSONRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*JSONConnection `json:"connections"`
	Devices     []*JSONDevice     `json:"devices"`
}

// NewJSONRobot returns a JSON representation of the robot.
func NewJSONRobot(robot *Robot) *JSONRobot {
	jsonRobot := &JSONRobot{
		Name:        robot.Name,
		Commands:    []string{},
		Connections: []*JSONConnection{},
		Devices:     []*JSONDevice{},
	}

	for command := range robot.Commands() {
		jsonRobot.Commands = append(jsonRobot.Commands, command)
	}

	robot.Devices().Each(func(device Device) {
		jsonDevice := NewJSONDevice(device)
		jsonRobot.Connections = append(jsonRobot.Connections, NewJSONConnection(robot.Connection(jsonDevice.Connection)))
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	})
	return jsonRobot
}

// Robot is a named entitity that manages a collection of connections and devices.
// It containes it's own work routine and a collection of
// custom commands to control a robot remotely via the Gobot api.
type Robot struct {
	Name        string
	Work        func()
	connections *connections
	devices     *devices
	Commander
	Eventer
}

type robots []*Robot

// Len counts the robots associated with this instance.
func (r *robots) Len() int {
	return len(*r)
}

// Start initialises the event loop. All robots that were added will
// be automtically started as a result of this call.
func (r *robots) Start() (errs []error) {
	for _, robot := range *r {
		if errs = robot.Start(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = errors.New(fmt.Sprintf("Robot %q: %v", robot.Name, err))
			}
			return
		}
	}
	return
}

// Each enumerates thru the robots and calls specified function
func (r *robots) Each(f func(*Robot)) {
	for _, robot := range *r {
		f(robot)
	}
}

// NewRobot returns a new Robot given a name and optionally accepts:
//
// 	[]Connection: Connections which are automatically started and stopped with the robot
//	[]Device: Devices which are automatically started and stopped with the robot
//	func(): The work routine the robot will execute once all devices and connections have been initialized and started
func NewRobot(name string, v ...interface{}) *Robot {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	r := &Robot{
		Name:        name,
		connections: &connections{},
		devices:     &devices{},
		Work:        nil,
		Eventer:     NewEventer(),
		Commander:   NewCommander(),
	}

	log.Println("Initializing Robot", r.Name, "...")

	for i := range v {
		switch v[i].(type) {
		case []Connection:
			log.Println("Initializing connections...")
			for _, connection := range v[i].([]Connection) {
				c := r.AddConnection(connection)
				log.Println("Initializing connection", c.Name(), "...")
			}
		case []Device:
			log.Println("Initializing devices...")
			for _, device := range v[i].([]Device) {
				d := r.AddDevice(device)
				log.Println("Initializing device", d.Name(), "...")
			}
		case func():
			r.Work = v[i].(func())
		default:
			log.Println("Unknown argument passed to NewRobot")
		}
	}

	return r
}

// Start a robot instance and runs it's work function if any. You should not
// need to manually start a robot if already part of a Gobot application as the
// robot will be automatically started for you.
func (r *Robot) Start() (errs []error) {
	log.Println("Starting Robot", r.Name, "...")
	if cerrs := r.Connections().Start(); len(cerrs) > 0 {
		errs = append(errs, cerrs...)
		return
	}
	if derrs := r.Devices().Start(); len(derrs) > 0 {
		errs = append(errs, derrs...)
		return
	}
	if r.Work != nil {
		log.Println("Starting work...")
		r.Work()
	}
	return
}

// Devices returns all devices associated with this robot.
func (r *Robot) Devices() *devices {
	return r.devices
}

// AddDevice adds a new device to the robots collection of devices. Returns the
// added device.
func (r *Robot) AddDevice(d Device) Device {
	*r.devices = append(*r.Devices(), d)
	return d
}

// Device returns a device given a name. Returns nil on no device.
func (r *Robot) Device(name string) Device {
	if r == nil {
		return nil
	}
	for _, device := range *r.devices {
		if device.Name() == name {
			return device
		}
	}
	return nil
}

// Connections returns all connections associated with this robot.
func (r *Robot) Connections() *connections {
	return r.connections
}

// AddConnection adds a new connection to the robots collection of connections.
// Returns the added connection.
func (r *Robot) AddConnection(c Connection) Connection {
	*r.connections = append(*r.Connections(), c)
	return c
}

// Connection returns a connection given a name. Returns nil on no connection.
func (r *Robot) Connection(name string) Connection {
	if r == nil {
		return nil
	}
	for _, connection := range *r.connections {
		if connection.Name() == name {
			return connection
		}
	}
	return nil
}
