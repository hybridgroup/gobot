package gobot

import (
  "time"
)

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