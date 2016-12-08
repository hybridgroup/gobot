package gobot_test

import (
	"fmt"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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

func ExampleRand() {
	i := gobot.Rand(100)
	fmt.Printf("%v is > 0 && < 100", i)
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
	gobottest.Assert(t, a, b)
}

func ExampleRefute() {
	t := &testing.T{}
	var a int = 100
	var b int = 200
	gobottest.Refute(t, a, b)
}
