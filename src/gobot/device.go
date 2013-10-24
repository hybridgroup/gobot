package gobot

import (
  "fmt"
  "strconv"
)

type Device struct {
  Name string
  Pin string
  Connection string
  Interval string
  Driver string
  Params map[string]string
  Robot Robot
}

type DeviceType struct {
  Name string
  Pin string
  Robot Robot
  Connection ConnectionType 
  Interval float64
  Driver Driver
  Params map[string]string
}

func NewDevice(d Device) *DeviceType {
  dt := new(DeviceType)
  dt.Name = d.Name
  dt.Pin = d.Pin
  dt.Robot = d.Robot
  dt.Params = d.Params
//  dt.Connection = determine_connection(params[:connection]) || default_connection
  dt.Connection = ConnectionType{Name: d.Connection,}
  if d.Interval == "" {
    dt.Interval = 0.5
  } else {
    f, err := strconv.ParseFloat(d.Interval, 64)
    if err == nil {
      dt.Interval = f
    } else {
      fmt.Println(err)
      dt.Interval = 0.5
    }
  }

  //dt.Driver = Driver.New(Driver{Robot: d.Robot, Params: d.Params,})
  return dt
}

func (dt *DeviceType) Start() {
  fmt.Println("Device " + dt.Name + "started")
  dt.Driver.Start()
}
    

//    def publish(event, *data)
//      if data.first
//        driver.publish(event_topic_name(event), *data)
//      else
//        driver.publish(event_topic_name(event))
//      end
//    end


// Execute driver command
func (dt *DeviceType) Command(method_name string, arguments []string) {
  //dt.Driver.Command(method_name, arguments)
}
