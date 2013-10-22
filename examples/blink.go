package main

import (
  . "gobot"
  "time"
  "fmt"
)

func main() {

  connections := []Connection{
    Connection {
      Name: "arduino",
      Adaptor: "arduino",
      Port: "/dev/ttyACM0",
    },
  }

  devices := []Device{
    Device{
      Name: "led",
      Driver: "arduino",
      Pin: "13",
    },
  }

  work := func(){
    Every(300 * time.Millisecond, func(){ fmt.Println("Greetings Human")})
  }
  
  robot := Robot{
      Connections: connections, 
      Devices: devices,
      Work: work,
  }

  robot.Start()
}
