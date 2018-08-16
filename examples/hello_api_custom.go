// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"net/http"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
)

func main() {
	master := gobot.NewMaster()

	a := api.NewAPI(master)

	// creates routes/handlers for the custom API
	a.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	})
	a.Get("/api/hello", func(res http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf("This command is attached to the robot %v", master.Robot("hello").Name)
		res.Write([]byte(msg))
	})

	// starts the API without the "standard" C2PIO API or Robeaux web interface.
	a.StartRaw()

	master.AddRobot(gobot.NewRobot("hello"))

	master.Start()
}
