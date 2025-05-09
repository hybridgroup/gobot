package system

import (
	"fmt"
	"sync"
	"time"
)

func startEdgePolling(
	pinLabel string,
	pinReadFunc func() (int, error),
	pollInterval time.Duration,
	wantedEdge int,
	eventHandler func(offset int, t time.Duration, et string, sn uint32, lsn uint32),
	quitChan chan struct{},
) error {
	if eventHandler == nil {
		return fmt.Errorf("an event handler is mandatory for edge polling")
	}
	if quitChan == nil {
		return fmt.Errorf("the quit channel is mandatory for edge polling")
	}

	const allEdges = "all"

	triggerEventOn := "none"
	switch wantedEdge {
	case digitalPinEventOnFallingEdge:
		triggerEventOn = DigitalPinEventFallingEdge
	case digitalPinEventOnRisingEdge:
		triggerEventOn = DigitalPinEventRisingEdge
	case digitalPinEventOnBothEdges:
		triggerEventOn = allEdges
	default:
		return fmt.Errorf("unsupported edge type %d for edge polling", wantedEdge)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		var oldState int
		var readStart time.Time
		var firstLoopDone bool
		for {
			select {
			case <-quitChan:
				if !firstLoopDone {
					wg.Done()
				}
				return
			default:
				// note: pure reading takes between 30us and 1ms on rasperry Pi1, typically 50us, with sysfs also 500us
				// can happen, so we use the time stamp before start of reading to reduce random duration offset
				readStart = time.Now()
				readValue, err := pinReadFunc()
				if err != nil {
					fmt.Printf("edge polling error occurred while reading the pin %s: %v\n", pinLabel, err)
					readValue = oldState // keep the value
				}
				if readValue != oldState {
					detectedEdge := DigitalPinEventRisingEdge
					if readValue < oldState {
						detectedEdge = DigitalPinEventFallingEdge
					}
					if firstLoopDone && (triggerEventOn == allEdges || triggerEventOn == detectedEdge) {
						eventHandler(0, time.Duration(readStart.UnixNano()), detectedEdge, 0, 0)
					}
					oldState = readValue
				}
				// the real poll interval is increased by the reading time, see also note above
				// negative or zero duration causes no sleep
				time.Sleep(pollInterval - time.Since(readStart))
				if !firstLoopDone {
					wg.Done()
					firstLoopDone = true
				}
			}
		}
	}()

	wg.Wait()
	return nil
}
