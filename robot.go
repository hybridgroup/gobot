package gobot

import (
	"fmt"
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
	connections   []*connection          `json:"-"`
	devices       []*device              `json:"-"`
	master        *Master                `json:"-"`
}

func (r *Robot) Start() {
	if r.master == nil {
		r.master = NewMaster()
	}

	r.master.Robots = []*Robot{r}
	r.master.Start()
}

func (r *Robot) startRobot() {
	r.initName()
	r.initCommands()
	r.initConnections()
	if r.startConnections() != true {
		panic("Could not start connections")
	}
	if r.initDevices() != true {
		panic("Could not initialize devices")
	}
	if r.startDevices() != true {
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
	r.connections = make([]*connection, len(r.Connections))
	log.Println("Initializing connections...")
	for i, connection := range r.Connections {
		log.Println("Initializing connection ", FieldByNamePtr(connection, "Name"), "...")
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

func (r *Robot) startConnections() bool {
	log.Println("Starting connections...")
	success := true
	for _, connection := range r.connections {
		log.Println("Starting connection " + connection.Name + "...")
		if connection.Connect() == false {
			success = false
			break
		}
	}
	return success
}

func (r *Robot) startDevices() bool {
	log.Println("Starting devices...")
	success := true
	for _, device := range r.devices {
		log.Println("Starting device " + device.Name + "...")
		if device.Start() == false {
			success = false
			break
		}
	}
	return success
}

func (r *Robot) haltDevices() {
	for _, device := range r.devices {
		device.Halt()
	}
}

func (r *Robot) finalizeConnections() {
	for _, connection := range r.connections {
		connection.Finalize()
	}
}

func (r *Robot) GetDevices() []*device {
	return r.devices
}

func (r *Robot) GetDevice(name string) *device {
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

func (r *Robot) GetConnections() []*connection {
	return r.connections
}

func (r *Robot) GetConnection(name string) *connection {
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
