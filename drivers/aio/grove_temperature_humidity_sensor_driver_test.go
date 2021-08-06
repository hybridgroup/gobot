package aio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

func TestGroveDHT11SensorDriver(t *testing.T) {
	// Arrange
	testAdaptor := newAioTestAdaptor()
	// Act
	dht11 := NewGroveDHT11SensorDriver(testAdaptor, "123")
	// Assert
	gobottest.Assert(t, dht11.Connection(), testAdaptor)
	gobottest.Assert(t, dht11.Pin(), "123")
	gobottest.Assert(t, dht11.interval, 1000*time.Millisecond)
}

func TestGroveDHT11SensorPublishTemperature(t *testing.T) {
	// Arrange
	sem := make(chan bool, 1)
	adaptor := newAioTestAdaptor()
	dht11 := NewGroveDHT11SensorDriver(adaptor, "1")
	expectedTemp := float32(25)
	expectedHum := float32(75)

	adaptor.TestAdaptorDHTRead(func() (t, h float32, err error) {
		return expectedTemp, expectedHum, nil
	})

	// Act
	dht11.Once(dht11.Event(Data), func(data interface{}) {
		gobottest.Assert(t, data.(GroveDHT11SensorState), GroveDHT11SensorState{
			Temperature: expectedTemp,
			Humidity:    expectedHum,
		})
		sem <- true
	})

	// Assert
	gobottest.Assert(t, dht11.Start(), nil)

	select {
	case <-sem:
	case <-time.After(2 * time.Second):
		t.Errorf(`Grove DHT11 Sensor Event "Data" was not published`)
	}

	gobottest.Assert(t, dht11.Temperature(), expectedTemp)
	gobottest.Assert(t, dht11.Humidity(), expectedHum)
}

func TestGroveDHT11SensorPublishError(t *testing.T) {
	// Arrange
	sem := make(chan bool, 1)
	adaptor := newAioTestAdaptor()
	dht11 := NewGroveDHT11SensorDriver(adaptor, "1")
	expectedErr := errors.New("failed to get data")

	adaptor.TestAdaptorDHTRead(func() (t, h float32, err error) {
		return t, h, expectedErr
	})

	// Act
	dht11.Once(dht11.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error), expectedErr)
		sem <- true
	})

	// Assert
	gobottest.Assert(t, dht11.Start(), nil)

	select {
	case <-sem:
	case <-time.After(2 * time.Second):
		t.Errorf(`Grove DHT11 Sensor Event "Error" was not published`)
	}
}

func TestGroveDHT11SensorHalt(t *testing.T) {
	// Arrange
	done := make(chan struct{})
	dht11 := NewGroveDHT11SensorDriver(newAioTestAdaptor(), "1")
	go func() {
		<-dht11.halt
		close(done)
	}()

	// Act
	gobottest.Assert(t, dht11.Halt(), nil)

	// Assert
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Grove DHT11 Sensorwas not halted")
	}
}

func TestGroveDHT11SensorDefaultName(t *testing.T) {
	// Arrange
	dht11 := NewGroveDHT11SensorDriver(newAioTestAdaptor(), "1")
	// Assert
	gobottest.Assert(t, strings.HasPrefix(dht11.Name(), "GroveDHT11Sensor"), true)
}

func TestGroveDHT11SensorSetName(t *testing.T) {
	// Arrange
	dht11 := NewGroveDHT11SensorDriver(newAioTestAdaptor(), "1")
	// Act
	dht11.SetName("mysensor")
	// Assert
	gobottest.Assert(t, dht11.Name(), "mysensor")
}

func TestGroveDHT11SensorUseInterval(t *testing.T) {
	// Arrange
	expectedInterval := time.Duration(5000 * time.Millisecond)
	dht11 := NewGroveDHT11SensorDriver(newAioTestAdaptor(), "1", WithGroveDHT11SensorInterval(expectedInterval))
	// Assert
	gobottest.Assert(t, dht11.interval, expectedInterval)
}

func TestGroveDHT11SensorUseInvalidInterval(t *testing.T) {
	// Arrange
	expectedInterval := time.Duration(1000 * time.Millisecond)
	dht11 := NewGroveDHT11SensorDriver(newAioTestAdaptor(), "1", WithGroveDHT11SensorInterval(300*time.Millisecond))
	// Assert
	gobottest.Assert(t, dht11.interval, expectedInterval)
}
