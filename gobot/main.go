package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gobot"
	app.Version = "0.0.1"
	app.Usage = "Command Line Utility for Gobot"
	app.Commands = []cli.Command{
		Generate(),
	}
	app.Run(os.Args)
}
