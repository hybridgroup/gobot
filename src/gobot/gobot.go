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

func Every(t time.Duration, f func()) {
  go func(){ 
    for{
      time.Sleep(t)  
      f()
    }
  }()
}

func After(t time.Duration, f func()) {
  go func(){ 
    time.Sleep(t)  
    f()
  }()
}