package gpio

import (
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

const (
	hcsr04SoundSpeed = 343 // in [m/s]
	// the device can measure 2 cm .. 4 m, this means sweep distances between 4 cm and 8 m
	// this cause pulse duration between 0.12 ms and 24 ms (at 34.3 cm/ms, ~0.03 ms/cm, ~3 ms/m)
	// so we use 60 ms as a limit for timeout and 100 ms for duration between 2 consecutive measurements
	hcsr04StartTransmitTimeout time.Duration = 100 * time.Millisecond // unfortunately takes sometimes longer than 60 ms
	hcsr04ReceiveTimeout       time.Duration = 60 * time.Millisecond
	hcsr04EmitTriggerDuration  time.Duration = 10 * time.Microsecond // according to specification
	hcsr04MonitorUpdate        time.Duration = 200 * time.Millisecond
	// the resolution of the device is ~3 mm, which relates to 10 us (343 mm/ms = 0.343 mm/us)
	// the poll interval increases the reading interval to this value and adds around 3 mm inaccuracy
	// it takes only an effect for fast systems, because reading inputs is typically much slower, e.g. 30-50 us on raspi
	// so, using the internal edge detection with "cdev" is more precise
	hcsr04PollInputIntervall time.Duration = 10 * time.Microsecond
)

// HCSR04Driver is a driver for ultrasonic range measurement.
type HCSR04Driver struct {
	*Driver
	triggerPinID                 string
	echoPinID                    string
	useEdgePolling               bool        // use discrete edge polling instead "cdev" from gpiod
	measureMutex                 *sync.Mutex // to ensure that only one measurement is done at a time
	triggerPin                   gobot.DigitalPinner
	echoPin                      gobot.DigitalPinner
	lastMeasureMicroSec          int64 // ~120 .. 24000 us
	distanceMonitorStopChan      chan struct{}
	distanceMonitorStopWaitGroup *sync.WaitGroup
	delayMicroSecChan            chan int64    // channel for event handler return value
	pollQuitChan                 chan struct{} // channel for quit the continuous polling
}

// NewHCSR04Driver creates a new instance of the driver for HC-SR04 (same as SEN-US01).
//
// Datasheet: https://www.makershop.de/download/HCSR04-datasheet-version-1.pdf
func NewHCSR04Driver(a gobot.Adaptor, triggerPinID string, echoPinID string, useEdgePolling bool) *HCSR04Driver {
	h := HCSR04Driver{
		Driver:         NewDriver(a, "HCSR04"),
		triggerPinID:   triggerPinID,
		echoPinID:      echoPinID,
		useEdgePolling: useEdgePolling,
		measureMutex:   &sync.Mutex{},
	}

	h.afterStart = func() error {
		tpin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(triggerPinID)
		if err != nil {
			return fmt.Errorf("error on get trigger pin: %v", err)
		}
		if err := tpin.ApplyOptions(system.WithPinDirectionOutput(0)); err != nil {
			return fmt.Errorf("error on apply output for trigger pin: %v", err)
		}
		h.triggerPin = tpin

		// pins are inputs by default
		epin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(echoPinID)
		if err != nil {
			return fmt.Errorf("error on get echo pin: %v", err)
		}

		epinOptions := []func(gobot.DigitalPinOptioner) bool{system.WithPinEventOnBothEdges(h.createEventHandler())}
		if h.useEdgePolling {
			h.pollQuitChan = make(chan struct{})
			epinOptions = append(epinOptions, system.WithPinPollForEdgeDetection(hcsr04PollInputIntervall, h.pollQuitChan))
		}
		if err := epin.ApplyOptions(epinOptions...); err != nil {
			return fmt.Errorf("error on apply options for echo pin: %v", err)
		}
		h.echoPin = epin

		h.delayMicroSecChan = make(chan int64)

		return nil
	}

	h.beforeHalt = func() error {
		if useEdgePolling {
			close(h.pollQuitChan)
		}

		if err := h.stopDistanceMonitor(); err != nil {
			fmt.Printf("no need to stop distance monitoring: %v\n", err)
		}

		// note: Unexport() of all pins will be done on adaptor.Finalize()

		close(h.delayMicroSecChan)

		return nil
	}

	return &h
}

// MeasureDistance retrieves the distance in front of sensor in meters and returns the measure. It is not designed
// to work in a fast loop! For this specific usage, use StartDistanceMonitor() associated with Distance() instead.
func (h *HCSR04Driver) MeasureDistance() (float64, error) {
	err := h.measureDistance()
	if err != nil {
		return 0, err
	}
	return h.Distance(), nil
}

// Distance returns the last distance measured in meter, it does not trigger a distance measurement
func (h *HCSR04Driver) Distance() float64 {
	distMm := h.lastMeasureMicroSec * hcsr04SoundSpeed / 1000 / 2
	return float64(distMm) / 1000.0
}

// StartDistanceMonitor starts continuous measurement. The current value can be read by Distance()
func (h *HCSR04Driver) StartDistanceMonitor() error {
	// ensure that start and stop can not interfere
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.distanceMonitorStopChan != nil {
		return fmt.Errorf("distance monitor already started for '%s'", h.name)
	}

	h.distanceMonitorStopChan = make(chan struct{})
	h.distanceMonitorStopWaitGroup = &sync.WaitGroup{}
	h.distanceMonitorStopWaitGroup.Add(1)

	go func(name string) {
		defer h.distanceMonitorStopWaitGroup.Done()
		for {
			select {
			case <-h.distanceMonitorStopChan:
				h.distanceMonitorStopChan = nil
				return
			default:
				if err := h.measureDistance(); err != nil {
					fmt.Printf("continuous measure distance skipped for '%s': %v\n", name, err)
				}
				time.Sleep(hcsr04MonitorUpdate)
			}
		}
	}(h.name)

	return nil
}

// StopDistanceMonitor stop the monitor process
func (h *HCSR04Driver) StopDistanceMonitor() error {
	// ensure that start and stop can not interfere
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.stopDistanceMonitor()
}

func (h *HCSR04Driver) createEventHandler() func(int, time.Duration, string, uint32, uint32) {
	var startTimestamp time.Duration
	return func(offset int, t time.Duration, et string, sn uint32, lsn uint32) {
		switch et {
		case system.DigitalPinEventRisingEdge:
			startTimestamp = t
		case system.DigitalPinEventFallingEdge:
			// unfortunately there is an additional falling edge at each start trigger, so we need to filter this
			// we use the start duration value for filtering
			if startTimestamp == 0 {
				return
			}
			h.delayMicroSecChan <- (t - startTimestamp).Microseconds()
			startTimestamp = 0
		}
	}
}

func (h *HCSR04Driver) stopDistanceMonitor() error {
	if h.distanceMonitorStopChan == nil {
		return fmt.Errorf("distance monitor is not yet started for '%s'", h.name)
	}

	h.distanceMonitorStopChan <- struct{}{}
	h.distanceMonitorStopWaitGroup.Wait()

	return nil
}

func (h *HCSR04Driver) measureDistance() error {
	h.measureMutex.Lock()
	defer h.measureMutex.Unlock()

	if err := h.emitTrigger(); err != nil {
		return err
	}

	// stop the loop if the measure is done or the timeout is elapsed
	timeout := hcsr04StartTransmitTimeout + hcsr04ReceiveTimeout
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout %s reached while waiting for value with echo pin %s", timeout, h.echoPinID)
	case h.lastMeasureMicroSec = <-h.delayMicroSecChan:
	}

	return nil
}

func (h *HCSR04Driver) emitTrigger() error {
	if err := h.triggerPin.Write(1); err != nil {
		return err
	}
	time.Sleep(hcsr04EmitTriggerDuration)
	return h.triggerPin.Write(0)
}
