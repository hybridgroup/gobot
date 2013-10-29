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

func On(cs chan interface{}) interface{}{
  for s := range cs {
    return s
  }
  return nil
}

func Work(robots []Robot) {
  for s := range robots {
    go robots[s].Start()
  }
  for{time.Sleep(10 * time.Millisecond)}
}
