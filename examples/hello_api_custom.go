//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"net/http"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
)

func main() {
	manager := gobot.NewManager()

	a := api.NewAPI(manager)

	// creates routes/handlers for the custom API
	a.Get("/", func(res http.ResponseWriter, req *http.Request) {
		if _, err := res.Write([]byte("OK")); err != nil {
			fmt.Println(err)
		}
	})
	a.Get("/api/hello", func(res http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf("This command is attached to the robot %v", manager.Robot("hello").Name)
		if _, err := res.Write([]byte(msg)); err != nil {
			fmt.Println(err)
		}
	})

	// starts the API without the default C2PIO API and Robeaux web interface.
	a.StartWithoutDefaults()

	manager.AddRobot(gobot.NewRobot("hello"))

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
