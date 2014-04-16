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

	driverTest, _ := template.New("").Parse(driverTest())
	file_location = fmt.Sprintf("%s/%s_driver_test.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	driverTest.Execute(f, name)
	f.Close()

	adaptorTest, _ := template.New("").Parse(adaptorTest())
	file_location = fmt.Sprintf("%s/%s_adaptor_test.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	adaptorTest.Execute(f, name)
	f.Close()

	testSuite, _ := template.New("").Parse(testSuite())
	file_location = fmt.Sprintf("%s/gobot-%s_suite_test.go", dir, args[0])
	fmt.Println("Creating", file_location)
	f, err = os.Create(file_location)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
	testSuite.Execute(f, name)
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

func (me *{{ .UpperName }}Adaptor) Disconnect() bool {
  return true
}

func (me *{{ .UpperName }}Adaptor) Finalize() bool {
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
	Adaptor *{{ .UpperName }}Adaptor
}

type {{ .UpperName }}Interface interface {
}

func New{{ .UpperName }}(adaptor *{{ .UpperName }}Adaptor) *{{ .UpperName }}Driver {
	d := new({{ .UpperName }}Driver)
	d.Events = make(map[string]chan interface{})
	d.Adaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *{{ .UpperName }}Driver) Init() bool { return true }
func (me *{{ .UpperName }}Driver) Start() bool { return true }
func (me *{{ .UpperName }}Driver) Halt() bool { return true }
`
}

func driverTest() string {
	return `package gobot{{ .UpperName }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("{{ .UpperName }}Driver", func() {
  var (
    driver *{{ .UpperName }}Driver
  )

  BeforeEach(func() {
    driver = New{{ .UpperName }}(new({{ .UpperName }}Adaptor))
  })

  It("Must be able to Start", func() {
    Expect(driver.Start()).To(Equal(true))
  })
  It("Must be able to Init", func() {
    Expect(driver.Init()).To(Equal(true))
  })
  It("Must be able to Halt", func() {
    Expect(driver.Halt()).To(Equal(true))
  })
})
`
}

func adaptorTest() string {
	return `package gobot{{ .UpperName }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("{{ .UpperName }}Adaptor", func() {
  var (
    adaptor *{{ .UpperName }}Adaptor
  )

  BeforeEach(func() {
    adaptor = new({{ .UpperName }}Adaptor)
  })

  It("Must be able to Finalize", func() {
    Expect(adaptor.Finalize()).To(Equal(true))
  })
  It("Must be able to Connect", func() {
    Expect(adaptor.Connect()).To(Equal(true))
  })
  It("Must be able to Disconnect", func() {
    Expect(adaptor.Disconnect()).To(Equal(true))
  })
  It("Must be able to Reconnect", func() {
    Expect(adaptor.Reconnect()).To(Equal(true))
  })
})
`
}

func testSuite() string {
	return `package gobot{{ .UpperName }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestGobot{{ .UpperName }}(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Gobot-{{ .UpperName }} Suite")
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
