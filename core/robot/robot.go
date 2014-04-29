package robot

import (
	"fmt"
	"github.com/hybridgroup/gobot/core/utils"
	"log"
	"math/rand"
	"time"
)

type Robot struct {
	Connections   []Connection           `json:"connections"`
	Devices       []Device               `json:"devices"`
	Name          string                 `json:"name"`
	Commands      map[string]interface{} `json:"-"`
	RobotCommands []string               `json:"commands"`
	Work          func()                 `json:"-"`
	connections   connections            `json:"-"`
	devices       devices                `json:"-"`
}

type Robots []*Robot

func (r Robots) Start() {
	for _, robot := range r {
		robot.Start()
	}
}

func (r Robots) Each(f func(*Robot)) {
	for _, robot := range r {
		f(robot)
	}
}

func NewRobot(name string, c []Connection, d []Device, work func()) *Robot {
	r := &Robot{
		Name:        name,
		Connections: c,
		Devices:     d,
		Work:        work,
	}
	r.initName()
	r.initCommands()
	r.initConnections()
	r.initDevices()
	return r
}

func (r *Robot) Start() {
	//	if !r.startConnections() {
	if err := r.GetConnections().Start(); err != nil {
		panic("Could not start connections")
	}
	if err := r.GetDevices().Start(); err != nil {
		panic("Could not start devices")
	}
	if r.Work != nil {
		r.Work()
	}
}

func (r *Robot) initName() {
	if r.Name == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		i := rand.Int()
		r.Name = fmt.Sprintf("Robot%v", i)
	}
}

func (r *Robot) initCommands() {
	r.RobotCommands = make([]string, 0)
	for k, _ := range r.Commands {
		r.RobotCommands = append(r.RobotCommands, k)
	}
}

func (r *Robot) initConnections() {
	r.connections = make(connections, len(r.Connections))
	log.Println("Initializing connections...")
	for i, connection := range r.Connections {
		log.Println("Initializing connection ", utils.FieldByNamePtr(connection, "Name"), "...")
		r.connections[i] = NewConnection(connection, r)
	}
}

func (r *Robot) initDevices() bool {
	r.devices = make([]*device, len(r.Devices))
	log.Println("Initializing devices...")
	for i, device := range r.Devices {
		r.devices[i] = NewDevice(device, r)
	}
	success := true
	for _, device := range r.devices {
		log.Println("Initializing device " + device.Name + "...")
		if device.Init() == false {
			success = false
			break
		}
	}
	return success
}

func (r *Robot) GetDevices() devices {
	return devices(r.devices)
}

func (r *Robot) Device(name string) *device {
	if r == nil {
		return nil
	}
	for _, device := range r.devices {
		if device.Name == name {
			return device
		}
	}
	return nil
}

func (r *Robot) GetConnections() connections {
	return connections(r.connections)
}

func (r *Robot) Connection(name string) *connection {
	if r == nil {
		return nil
	}
	for _, connection := range r.connections {
		if connection.Name == name {
			return connection
		}
	}
	return nil
}
