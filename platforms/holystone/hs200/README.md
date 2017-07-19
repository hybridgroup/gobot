## How to Use

```go
package main

import (
	"log"
	"time"

	"gobot.io/x/gobot/platforms/holystone/hs200"
	"fmt"
)

func main() {
	drone, err := hs200.NewDriver("172.16.10.1:8888", "172.16.10.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enable!")
	drone.Enable()
	time.Sleep(5 * time.Second)
	fmt.Println("Take off!")
	drone.TakeOff()
	time.Sleep(5 * time.Second)
	fmt.Println("Land!")
	drone.Land()
	time.Sleep(5 * time.Second)
	fmt.Println("Disable!")
	drone.Disable()
	time.Sleep(5*time.Second)
}
```