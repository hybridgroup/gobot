/*
Package api provides a webserver to interact with your Gobot program over the network.

Example:

    package main

    import (
    	"fmt"

    	"gobot.io/x/gobot"
    	"gobot.io/x/gobot/api"
    )

    func main() {
    	gbot := gobot.NewMaster()

      // Starts the API server on default port 3000
    	api.NewAPI(gbot).Start()

      // Accessible via http://localhost:3000/api/commands/say_hello
    	gbot.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return "Master says hello!"
    	})

    	hello := gbot.AddRobot(gobot.NewRobot("Eve"))

      // Accessible via http://localhost:3000/robots/Eve/commands/say_hello
    	hello.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return fmt.Sprintf("%v says hello!", hello.Name)
    	})

    	gbot.Start()
    }

It follows Common Protocol for Programming Physical Input and Output (CPPP-IO) spec:
https://gobot.io/x/cppp-io
*/
package api
