package gobot

import (
  "time"
)

//type Gobot struct{
//  Robot
//  Connections []Connection
//  Devices []Device
//  Name string
//  Work func()
//}

func Every(t time.Duration, ret func()) {
  go func(){ 
    for{
      ret()
      time.Sleep(t)  
    }
  }()
}