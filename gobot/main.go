package main

import (
	"fmt"
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

var g_cmd *commander.Command

func init() {
	g_cmd = &commander.Command{
		UsageLine: "gobot <command>",
		Subcommands: []*commander.Command{
			generate(),
		},
		Flag: *flag.NewFlagSet("gobot", flag.ExitOnError),
	}
}

func main() {
	err := g_cmd.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Printf("**err**: %v\n", err)
		os.Exit(1)
	}
	return
}
