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

// hcsr04OptionApplier needs to be implemented by each configurable option type
type hcsr04OptionApplier interface {
	apply(cfg *hcsr04Configuration)
}

// hcsr04Configuration contains all changeable attributes of the driver.
type hcsr04Configuration struct {
	useEdgePolling bool
}

// hcsr04UseEdgePollingOption is the type for applying to use discrete edge polling instead pin edge detection
// by "cdev" from the go-gpiocdev package.
type hcsr04UseEdgePollingOption bool

// HCSR04Driver is a driver for ultrasonic range measurement.
type HCSR04Driver struct {
	*driver
	hcsr04Cfg                    *hcsr04Configuration
	triggerPinID                 string
	echoPinID                    string
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
//
// Supported options:
//
//	"WithName"
func NewHCSR04Driver(a gobot.Adaptor, triggerPinID, echoPinID string, opts ...interface{}) *HCSR04Driver {
	d := HCSR04Driver{
		driver:       newDriver(a, "HCSR04"),
		hcsr04Cfg:    &hcsr04Configuration{},
		triggerPinID: triggerPinID,
		echoPinID:    echoPinID,
		measureMutex: &sync.Mutex{},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case hcsr04OptionApplier:
			o.apply(d.hcsr04Cfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	d.afterStart = func() error {
		tpin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(triggerPinID) //nolint:forcetypeassert // ok here
		if err != nil {
			return fmt.Errorf("error on get trigger pin: %v", err)
		}
		if err := tpin.ApplyOptions(system.WithPinDirectionOutput(0)); err != nil {
			return fmt.Errorf("error on apply output for trigger pin: %v", err)
		}
		d.triggerPin = tpin

		// pins are inputs by default
		epin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(echoPinID) //nolint:forcetypeassert // ok here
		if err != nil {
			return fmt.Errorf("error on get echo pin: %v", err)
		}

		epinOptions := []func(gobot.DigitalPinOptioner) bool{system.WithPinEventOnBothEdges(d.createEventHandler())}
		if d.hcsr04Cfg.useEdgePolling {
			d.pollQuitChan = make(chan struct{})
			epinOptions = append(epinOptions, system.WithPinPollForEdgeDetection(hcsr04PollInputIntervall, d.pollQuitChan))
		}
		if err := epin.ApplyOptions(epinOptions...); err != nil {
			return fmt.Errorf("error on apply options for echo pin: %v", err)
		}
		d.echoPin = epin

		d.delayMicroSecChan = make(chan int64)

		return nil
	}

	d.beforeHalt = func() error {
		if d.hcsr04Cfg.useEdgePolling {
			close(d.pollQuitChan)
		}

		if err := d.stopDistanceMonitor(); err != nil {
			fmt.Printf("no need to stop distance monitoring: %v\n", err)
		}

		// note: Unexport() of all pins will be done on adaptor.Finalize()

		close(d.delayMicroSecChan)

		return nil
	}

	return &d
}

// WithHCSR04UseEdgePolling use discrete edge polling instead pin edge detection by "cdev" from the go-gpiocdev package.
func WithHCSR04UseEdgePolling() hcsr04OptionApplier {
	return hcsr04UseEdgePollingOption(true)
}

// MeasureDistance retrieves the distance in front of sensor in meters and returns the measure. It is not designed
// to work in a fast loop! For this specific usage, use StartDistanceMonitor() associated with Distance() instead.
func (d *HCSR04Driver) MeasureDistance() (float64, error) {
	err := d.measureDistance()
	if err != nil {
		return 0, err
	}
	return d.Distance(), nil
}

// Distance returns the last distance measured in meter, it does not trigger a distance measurement
func (d *HCSR04Driver) Distance() float64 {
	distMm := d.lastMeasureMicroSec * hcsr04SoundSpeed / 1000 / 2
	return float64(distMm) / 1000.0
}

// StartDistanceMonitor starts continuous measurement. The current value can be read by Distance()
func (d *HCSR04Driver) StartDistanceMonitor() error {
	// ensure that start and stop can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.distanceMonitorStopChan != nil {
		return fmt.Errorf("distance monitor already started for '%s'", d.driverCfg.name)
	}

	d.distanceMonitorStopChan = make(chan struct{})
	d.distanceMonitorStopWaitGroup = &sync.WaitGroup{}
	d.distanceMonitorStopWaitGroup.Add(1)

	go func(name string) {
		defer d.distanceMonitorStopWaitGroup.Done()
		for {
			select {
			case <-d.distanceMonitorStopChan:
				d.distanceMonitorStopChan = nil
				return
			default:
				if err := d.measureDistance(); err != nil {
					fmt.Printf("continuous measure distance skipped for '%s': %v\n", name, err)
				}
				time.Sleep(hcsr04MonitorUpdate)
			}
		}
	}(d.driverCfg.name)

	return nil
}

// StopDistanceMonitor stop the monitor process
func (d *HCSR04Driver) StopDistanceMonitor() error {
	// ensure that start and stop can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.stopDistanceMonitor()
}

func (d *HCSR04Driver) createEventHandler() func(int, time.Duration, string, uint32, uint32) {
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
			d.delayMicroSecChan <- (t - startTimestamp).Microseconds()
			startTimestamp = 0
		}
	}
}

func (d *HCSR04Driver) stopDistanceMonitor() error {
	if d.distanceMonitorStopChan == nil {
		return fmt.Errorf("distance monitor is not yet started for '%s'", d.driverCfg.name)
	}

	d.distanceMonitorStopChan <- struct{}{}
	d.distanceMonitorStopWaitGroup.Wait()

	return nil
}

func (d *HCSR04Driver) measureDistance() error {
	d.measureMutex.Lock()
	defer d.measureMutex.Unlock()

	if err := d.emitTrigger(); err != nil {
		return err
	}

	// stop the loop if the measure is done or the timeout is elapsed
	timeout := hcsr04StartTransmitTimeout + hcsr04ReceiveTimeout
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout %s reached while waiting for value with echo pin %s", timeout, d.echoPinID)
	case d.lastMeasureMicroSec = <-d.delayMicroSecChan:
	}

	return nil
}

func (d *HCSR04Driver) emitTrigger() error {
	if err := d.triggerPin.Write(1); err != nil {
		return err
	}
	time.Sleep(hcsr04EmitTriggerDuration)
	return d.triggerPin.Write(0)
}

func (o hcsr04UseEdgePollingOption) String() string {
	return "hcsr04 use edge polling option"
}

func (o hcsr04UseEdgePollingOption) apply(cfg *hcsr04Configuration) {
	cfg.useEdgePolling = bool(o)
}
