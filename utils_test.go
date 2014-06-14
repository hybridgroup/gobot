package gobot

import (
	"fmt"
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	i := 0
	Every(2*time.Millisecond, func() {
		i++
	})
	time.Sleep(5 * time.Millisecond)
	Expect(t, i, 2)
}

func TestAfter(t *testing.T) {
	i := 0
	After(1*time.Millisecond, func() {
		i++
	})
	time.Sleep(2 * time.Millisecond)
	Expect(t, i, 1)
}

func TestPublish(t *testing.T) {
	e := &Event{Chan: make(chan interface{}, 1)}
	Publish(e, 1)
	Publish(e, 2)
	Publish(e, 3)
	Publish(e, 4)
	i := <-e.Chan
	Expect(t, i, 1)
}

func TestOn(t *testing.T) {
	var i int
	e := NewEvent()
	On(e, func(data interface{}) {
		i = data.(int)
	})
	Publish(e, 10)
	time.Sleep(1 * time.Millisecond)
	Expect(t, i, 10)
}

func TestFromScale(t *testing.T) {
	Expect(t, FromScale(5, 0, 10), 0.5)
}

func TestToScale(t *testing.T) {
	Expect(t, ToScale(500, 0, 10), 10.0)
	Expect(t, ToScale(-1, 0, 10), 0.0)
	Expect(t, ToScale(0.5, 0, 10), 5.0)
}

func TestRand(t *testing.T) {
	a := Rand(1000)
	b := Rand(1000)
	if a == b {
		t.Error(fmt.Sprintf("%v should not equal %v", a, b))
	}
}

func TestFieldByName(t *testing.T) {
	testInterface := *NewTestStruct()
	Expect(t, FieldByName(testInterface, "i").Int(), int64(10))
}

func TestFieldByNamePtr(t *testing.T) {
	testInterface := NewTestStruct()
	Expect(t, FieldByNamePtr(testInterface, "f").Float(), 0.2)
}

func TestCall(t *testing.T) {
	testInterface := NewTestStruct()
	Expect(t, Call(testInterface, "Hello", "Human", "How are you?")[0].String(), "Hello Human! How are you?")
}
