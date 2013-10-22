package main

import (
  . "gobot"
  "time"
  "fmt"
)

func main() {
  
  robot := Robot{
    Work: func(){
      Every(300 * time.Millisecond, func(){ fmt.Println("Greetings human") })
    },
  }

  robot.Start()
}
