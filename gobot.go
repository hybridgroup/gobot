package gobot

import (
  "time"
  "math/rand"
  "net"
)

func Every(t string, f func()) {
  dur := parseDuration(t)
  go func(){ 
    for{
      time.Sleep(dur)
      go f()
    }
  }()
}

func After(t string, f func()) {
  dur := parseDuration(t)
  go func(){ 
    time.Sleep(dur)
    f()
  }()
}

func parseDuration(t string) time.Duration {
  return ParseDuration(t)
}
func ParseDuration(t string) time.Duration {
  dur, err := time.ParseDuration(t)
  if err != nil {
    panic(err)
  }
  return dur
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

func ConnectTo(port string) net.Conn {
 tcpPort, err := net.Dial("tcp", port)
 if err != nil {
  panic(err)
 }
 return tcpPort
}
