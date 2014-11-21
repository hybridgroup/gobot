package gobot_test

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"testing"
	"time"
)

func ExampleEvery() {
	gobot.Every(1*time.Second, func() {
		fmt.Println("Hello")
	})
}

func ExampleAfter() {
	gobot.After(1*time.Second, func() {
		fmt.Println("Hello")
	})
}

func ExamplePublish() {
	e := gobot.NewEvent()
	gobot.Publish(e, 100)
}

func ExampleOn() {
	e := gobot.NewEvent()
	gobot.On(e, func(s interface{}) {
		fmt.Println(s)
	})
	gobot.Publish(e, 100)
	gobot.Publish(e, 200)
}

func ExampleOnce() {
	e := gobot.NewEvent()
	gobot.Once(e, func(s interface{}) {
		fmt.Println(s)
		fmt.Println("I will no longer respond to events")
	})
	gobot.Publish(e, 100)
	gobot.Publish(e, 200)
}

func ExampleRand() {
	i := gobot.Rand(100)
	fmt.Sprintln("%v is > 0 && < 100", i)
}

func ExampleFromScale() {
	fmt.Println(gobot.FromScale(5, 0, 10))
	// Output:
	// 0.5
}

func ExampleToScale() {
	fmt.Println(gobot.ToScale(500, 0, 10))
	// Output:
	// 10
}

func ExampleAssert() {
	t := &testing.T{}
	var a int = 100
	var b int = 100
	gobot.Assert(t, a, b)
}

func ExampleRefute() {
	t := &testing.T{}
	var a int = 100
	var b int = 200
	gobot.Refute(t, a, b)
}
