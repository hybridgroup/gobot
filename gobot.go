package gobot

import (
  "time"
  "math/rand"
)

func Every(t string, f func()) {
  dur,_ := time.ParseDuration(t)
  go func(){ 
    for{
      time.Sleep(dur)
      f()
    }
  }()
}

func After(t string, f func()) {
  dur,_ := time.ParseDuration(t)
  go func(){ 
    time.Sleep(dur)
    f()
  }()
}

func Random(min int, max int) int {
  rand.Seed(time.Now().UTC().UnixNano())
  return rand.Intn(max - min) + min
}