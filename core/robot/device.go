package robot

import (
	"errors"
	"github.com/hybridgroup/gobot/core/driver"
	"github.com/hybridgroup/gobot/core/utils"
	"log"
	"reflect"
)

type Device interface {
	Init() bool
	Start() bool
	Halt() bool
}

type device struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"driver"`
	Interval string                 `json:"-"`
	Robot    *Robot                 `json:"-"`
	Driver   driver.DriverInterface `json:"-"`
}

type devices []*device

// Start() starts all the devices.
func (d devices) Start() error {
	var err error
	log.Println("Starting devices...")
	for _, device := range d {
		log.Println("Starting device " + device.Name + "...")
		if device.Start() == false {
			err = errors.New("Could not start connection")
			break
		}
	}
	return err
}

// Halt() stop all the devices.
func (d devices) Halt() {
	for _, device := range d {
		device.Halt()
	}
}

func NewDevice(driver driver.DriverInterface, r *Robot) *device {
	d := new(device)
	s := reflect.ValueOf(driver).Type().String()
	d.Type = s[1:len(s)]
	d.Name = utils.FieldByNamePtr(driver, "Name").String()
	d.Robot = r
	if utils.FieldByNamePtr(driver, "Interval").String() == "" {
		utils.FieldByNamePtr(driver, "Interval").SetString("0.1s")
	}
	d.Driver = driver
	return d
}

func (d *device) Init() bool {
	log.Println("Device " + d.Name + " initialized")
	return d.Driver.Init()
}

func (d *device) Start() bool {
	log.Println("Device " + d.Name + " started")
	return d.Driver.Start()
}

func (d *device) Halt() bool {
	log.Println("Device " + d.Name + " halted")
	return d.Driver.Halt()
}

func (d *device) Commands() interface{} {
	return utils.FieldByNamePtr(d.Driver, "Commands").Interface()
}
