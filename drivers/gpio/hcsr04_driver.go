package hcsr04

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	soundSpeed       float32       = 343.0
	measurementCycle time.Duration = 60 // 60ms between two measurements
	// MonitorUpdate is the time between each monitor update
	MonitorUpdate time.Duration = 100 * time.Millisecond
)

// HCSR04 instance
type HCSR04 struct {
	name       string
	connection gobot.Adaptor
	triggerPin *gpio.DirectPinDriver
	echoPin    *gpio.DirectPinDriver
	// triggerPin             rpio.Pin
	// echoPin                rpio.Pin
	mux                    sync.Mutex
	Measure                float32 // The last measure
	distanceMonitorControl chan int
	distanceMonitorStarted bool
	gobot.Commander
}

// NewHCSR04 creates a new HCSR04 instance
func NewHCSR04(a gobot.Adaptor, triggerPinID string, echoPinID string) *HCSR04 {
	hcsr04 := &HCSR04{
		name:       gobot.DefaultName("HCSR04"),
		triggerPin: gpio.NewDirectPinDriver(a, triggerPinID),
		echoPin:    gpio.NewDirectPinDriver(a, echoPinID),
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	// hcsr04 := HCSR04{
	// 	triggerPinID: triggerPinID,
	// 	echoPinID:    echoPinID,
	// }
	// hcsr04.triggerPin = rpio.Pin(hcsr04.triggerPinID)
	// hcsr04.triggerPin.Mode(rpio.Output)

	// hcsr04.echoPin = rpio.Pin(hcsr04.echoPinID)
	// hcsr04.echoPin.Mode(rpio.Input)
	// hcsr04.echoPin.PullDown()
	// hcsr04.triggerPin.Low()
	return hcsr04
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
	readedValue, _ := hcsr04.echoPin.DigitalRead()
	if readedValue == 1 {
		return 0, errors.New("already receiving echo")
	}
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
		// readedValue, _ := pin.DigitalRead()
		readedValue := 1
		// fmt.Println("Lectura: ", readedValue, state, *quit)
		if readedValue != state && !*quit {
			break
		}
	}
	time.Sleep(100)
	readedValue, _ := pin.DigitalRead()
	if readedValue == state && !*quit {
		outChan <- time.Now().UnixNano()
	}
}

func pulseToDistance(pulseDuration int64) float32 {
	return float32(pulseDuration) / 1000000000.0 * soundSpeed / 2
}

// // GetDistance returns the last distance measured
// // Contrary to MeasureDistance, GetDistance does not trigger a distance measurement
func (hcsr04 *HCSR04) GetDistance() float32 {
	return hcsr04.Measure
}

// // StartDistanceMonitor starts a process which will keep Measure updated
func (hcsr04 *HCSR04) StartDistanceMonitor() error {
	hcsr04.distanceMonitorControl = make(chan int)
	if hcsr04.distanceMonitorStarted {
		return errors.New("monitor already started")
	}
	go hcsr04.distanceMonitor()
	return nil
}

// // StopDistanceMonitor stop the monitor process
func (hcsr04 *HCSR04) StopDistanceMonitor() {
	if hcsr04.distanceMonitorStarted {
		hcsr04.distanceMonitorControl <- 1
	}
}

func (hcsr04 *HCSR04) distanceMonitor() {
	for {
		select {
		case <-hcsr04.distanceMonitorControl:
			hcsr04.distanceMonitorStarted = false
			return
		default:
			if _, err := hcsr04.MeasureDistance(); err != nil {
				log.WithField("error", err).Error("impossible to measure distance")
			}
		}
		time.Sleep(MonitorUpdate)
	}
}

func (h *HCSR04) Name() string { return h.name }

func (h *HCSR04) SetName(n string) { h.name = n }

func (h *HCSR04) Start() (err error) { return }

func (h *HCSR04) Halt() (err error) { return }

func (h *HCSR04) Connection() gobot.Connection {
	return h.connection.(gobot.Connection)
}
