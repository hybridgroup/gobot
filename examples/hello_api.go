//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"html"
	"net/http"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
)

func main() {
	master := gobot.NewMaster()

	a := api.NewAPI(master)
	a.AddHandler(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))
	})
	a.Debug()
	a.Start()

	master.AddCommand("custom_gobot_command",
		func(params map[string]interface{}) interface{} {
			return "This command is attached to the mcp!"
		})

	hello := master.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	if err := master.Start(); err != nil {
		panic(err)
	}
}
