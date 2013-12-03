package main

import (
	"fmt"
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"os"
	"text/template"
	"unicode"
)

func generate() *commander.Command {
	cmd := &commander.Command{
		Run:       doGenerate,
		UsageLine: "generate [options]",
		Short:     "Generates a gobot library skeleton",
		Long: `
Generates a gobot library skeleton.

ex:
 $ gobot generate myProject
`,
		Flag: *flag.NewFlagSet("gobot-generate", flag.ExitOnError),
	}
	return cmd
}

type Generate struct {
	Name string
}

func doGenerate(cmd *commander.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(cmd.Long)
		return
	}
	pwd, _ := os.Getwd()
	dir := fmt.Sprintf("%s/gobot-%s", pwd, args[0])
	fmt.Println("Creating", dir)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		fmt.Println(err)
		err = nil
	}

	a := []rune(args[0])
	a[0] = unicode.ToUpper(a[0])
	s := string(a)

	name := Generate{Name: s}

	adaptor, _ := template.New("").Parse(adaptor())
	file_location := fmt.Sprintf("%s/%s_adaptor.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err := os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	adaptor.Execute(f, name)
	f.Close()

	file_location = fmt.Sprintf("%s/%s_driver.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	driver, _ := template.New("").Parse(driver())
	driver.Execute(f, name)
	f.Close()
}

func adaptor() string {
	return `package gobot{{ .Name }}

import (
	"github.com/hybridgroup/gobot"
)

type {{ .Name }}Adaptor struct {
	gobot.Adaptor
}

func (me *{{ .Name }}Adaptor) Connect() {
}

func (me *{{ .Name }}Adaptor) Disconnect() {
}
`
}

func driver() string {
	return `package gobot{{ .Name }}

import (
	"github.com/hybridgroup/gobot"
)

type {{ .Name }}Driver struct {
	gobot.Driver
	{{ .Name }}Adaptor *{{ .Name }}Adaptor
}

func New{{ .Name }}(adaptor *{{ .Name }}Adaptor) *{{ .Name }}Driver {
	d := new({{ .Name }}Driver)
	d.Events = make(map[string]chan interface{})
	d.{{ .Name }}Adaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *{{ .Name }}Driver) StartDriver() {
	gobot.Every(sd.Interval, func() {
		me.handleMessageEvents()
	})
}

func (me *{{ .Name }}Driver) handleMessageEvents() {
}
`
}
