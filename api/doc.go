/*
Package api provides functionally to expose your gobot programs
to other by using starting a web server and adding commands.

Example:

    package main

    import (
    	"fmt"

    	"github.com/hybridgroup/gobot"
    	"github.com/hybridgroup/gobot/api"
    )

    func main() {
    	gbot := gobot.NewGobot()

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
https://github.com/hybridgroup/cppp-io
*/
package api
