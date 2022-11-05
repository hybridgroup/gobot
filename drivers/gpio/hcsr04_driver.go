package gpio

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	soundSpeed float32 = 343.0 // in [m/s]
	// the device can measure 2cm .. 4m, this means sweep distances between 4cm and 8m
	// this cause pulse durations between 0.12ms and  24ms (at 34.3 cm/ms, ~0.03 ms/cm, ~3ms/m)
	measurementCycle time.Duration = 60 // 60ms between two measurements
	MonitorUpdate    time.Duration = 100 * time.Millisecond
)

// HCSR04 instance
type HCSR04 struct {
	name                   string
	connection             gobot.Adaptor
	triggerPin             *gpio.DirectPinDriver
	echoPin                *gpio.DirectPinDriver
	mux                    sync.Mutex
	Measure                float32 // The last measure
	distanceMonitorControl chan int
	distanceMonitorStarted bool
	gobot.Commander
}

// NewHCSR04 creates a new HCSR04 instance
func NewHCSR04(a gobot.Adaptor, triggerPinID string, echoPinID string) *HCSR04 {
	return &HCSR04{
		name:       gobot.DefaultName("HCSR04"),
		triggerPin: gpio.NewDirectPinDriver(a, triggerPinID),
		echoPin:    gpio.NewDirectPinDriver(a, echoPinID),
		connection: a,
		Commander:  gobot.NewCommander(),
	}
}

// MeasureDistance measure the distance in front of sensor in meters
// and returns the measure
// MeasureDistance triggers a distance measurement by the sensor
//
// ! MeasureDistance is not design to work in a fast loop
// For this specific usage, use StartDistanceMonitor associated with GetDistance Instead
func (hcsr04 *HCSR04) MeasureDistance() (float32, error) {
	hcsr04.mux.Lock()
	defer hcsr04.mux.Unlock()
	pulseDuration, err := hcsr04.measurePulse()
	if err != nil {
		return 0, err
	}
	hcsr04.Measure = pulseToDistance(pulseDuration)
	return hcsr04.Measure, nil
}

// GetDistance returns the last distance measured
// Contrary to MeasureDistance, GetDistance does not trigger a distance measurement
func (hcsr04 *HCSR04) GetDistance() float32 {
	return hcsr04.Measure
}

// StartDistanceMonitor starts a process which will keep Measure updated
func (hcsr04 *HCSR04) StartDistanceMonitor() error {
	hcsr04.distanceMonitorControl = make(chan int)
	if hcsr04.distanceMonitorStarted {
		return errors.New("monitor already started")
	}
	go hcsr04.distanceMonitor()
	return nil
}

// StopDistanceMonitor stop the monitor process
func (hcsr04 *HCSR04) StopDistanceMonitor() {
	if hcsr04.distanceMonitorStarted {
		hcsr04.distanceMonitorControl <- 1
	}
}

func (h *HCSR04) Name() string { return h.name }

func (h *HCSR04) SetName(n string) { h.name = n }

func (h *HCSR04) Start() (err error) { return }

func (h *HCSR04) Halt() (err error) { return }

func (h *HCSR04) Connection() gobot.Connection {
	return h.connection.(gobot.Connection)
}

func (hcsr04 *HCSR04) distanceMonitor() {
	for {
		select {
		case <-hcsr04.distanceMonitorControl:
			hcsr04.distanceMonitorStarted = false
			return
		default:
			if _, err := hcsr04.MeasureDistance(); err != nil {
				fmt.Println("error: impossible to measure distance", err)
			}
		}
		time.Sleep(MonitorUpdate)
	}
}

func (hcsr04 *HCSR04) emitTrigger() {
	hcsr04.triggerPin.On()
	time.Sleep(10 * time.Microsecond)
	hcsr04.triggerPin.Off()
}

func (hcsr04 *HCSR04) measurePulse() (int64, error) {
	startChan := make(chan int64)
	stopChan := make(chan int64)
	startQuit := false
	stopQuit := false
	var startTime int64
	var stopTime int64
	go getPinStateChangeTime(hcsr04.echoPin, 1, startChan, &startQuit)
	hcsr04.emitTrigger()
	select {
	case t := <-startChan:
		startTime = t
	case <-time.After(measurementCycle * time.Millisecond):
		startQuit = true
		return 0, fmt.Errorf("echo not received after %d milliseconds", measurementCycle)
	}
	go getPinStateChangeTime(hcsr04.echoPin, 0, stopChan, &stopQuit)
	select {
	case t := <-stopChan:
		stopTime = t
	case <-time.After(measurementCycle * time.Millisecond):
		stopQuit = true
		return 0, fmt.Errorf("echo received for more than %d milliseconds", measurementCycle)
	}
	return stopTime - startTime, nil
}

func getPinStateChangeTime(pin *gpio.DirectPinDriver, state int, outChan chan int64, quit *bool) {

	for {
		readedValue, _ := pin.DigitalRead()
		// stop the loop if the state is different or a quit is done
		if readedValue == state || *quit {
			break
		}
	}
	time.Sleep(100)
	readedValue, err := pin.DigitalRead()
	if err != nil {
		fmt.Println(err)
	}
	if readedValue == state && !*quit {
		outChan <- time.Now().UnixNano()
	}
}

func pulseToDistance(pulseDuration int64) float32 {
	return float32(pulseDuration) / 1000000000.0 * soundSpeed / 2
}
