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
		Short:     "Generates a Gobot library skeleton",
		Long: `
Generates a Gobot library skeleton.

ex:
 $ gobot generate myProject
`,
		Flag: *flag.NewFlagSet("gobot-generate", flag.ExitOnError),
	}
	return cmd
}

type Generate struct {
	Name      string
	UpperName string
}

func doGenerate(cmd *commander.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println(cmd.Long)
		return nil
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

	name := Generate{UpperName: s, Name: string(args[0])}

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

	driver, _ := template.New("").Parse(driver())
	file_location = fmt.Sprintf("%s/%s_driver.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	driver.Execute(f, name)
	f.Close()

	readme, _ := template.New("").Parse(readme())
	file_location = fmt.Sprintf("%s/README.md", dir)
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	readme.Execute(f, name)
	f.Close()

	file_location = fmt.Sprintf("%s/LICENSE", dir)
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	f.Close()
	return nil
}

func adaptor() string {
	return `package gobot{{ .UpperName }}

import (
	"github.com/hybridgroup/gobot"
)

type {{ .UpperName }}Adaptor struct {
	gobot.Adaptor
}

func (me *{{ .UpperName }}Adaptor) Connect() bool {
  return true
}

func (me *{{ .UpperName }}Adaptor) Reconnect() bool {
  return true
}

func (me *{{ .Name }}Adaptor) Disconnect() bool {
  return true
}

func (me *{{ .Name }}Adaptor) Finalize() bool {
  return true
}
`
}

func driver() string {
	return `package gobot{{ .UpperName }}

import (
	"github.com/hybridgroup/gobot"
)

type {{ .UpperName }}Driver struct {
	gobot.Driver
	{{ .UpperName }}Adaptor *{{ .UpperName }}Adaptor
}

type {{ .UpperName}}Interface interface {
}

func New{{ .UpperName }}(adaptor *{{ .UpperName }}Adaptor) *{{ .UpperName }}Driver {
	d := new({{ .UpperName }}Driver)
	d.Events = make(map[string]chan interface{})
	d.{{ .UpperName }}Adaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *{{ .UpperName }}Driver) Start() bool {
	gobot.Every(sd.Interval, func() {
		me.handleMessageEvents()
	})
  return true
}

func (me *{{ .UpperName }}Driver) handleMessageEvents() {
}
`
}

func readme() string {
	return `# Gobot for {{ .Name }}

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for {{ .Name }}.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Installing

    go get path/to/repo/gobot-{{ .Name }}

## Using

    your example code here...

## Connecting

Explain how to connect from the computer to the device here...

## License

Copyright (c) 2014 Your Name Here. See LICENSE for more details
`
}
