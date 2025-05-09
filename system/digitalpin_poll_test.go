package system

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_startEdgePolling(t *testing.T) {
	type readValue struct {
		value int
		err   string
	}
	tests := map[string]struct {
		eventOnEdge            int
		simulateReadValues     []readValue
		simulateNoEventHandler bool
		simulateNoQuitChan     bool
		wantEdgeTypes          []string
		wantErr                string
	}{
		"edge_falling": {
			eventOnEdge: digitalPinEventOnFallingEdge,
			simulateReadValues: []readValue{
				{value: 1},
				{value: 0},
				{value: 1},
				{value: 0},
				{value: 0},
			},
			wantEdgeTypes: []string{DigitalPinEventFallingEdge, DigitalPinEventFallingEdge},
		},
		"no_edge_falling": {
			eventOnEdge: digitalPinEventOnFallingEdge,
			simulateReadValues: []readValue{
				{value: 0},
				{value: 1},
				{value: 1},
			},
			wantEdgeTypes: nil,
		},
		"edge_rising": {
			eventOnEdge: digitalPinEventOnRisingEdge,
			simulateReadValues: []readValue{
				{value: 0},
				{value: 1},
				{value: 0},
				{value: 1},
				{value: 1},
			},
			wantEdgeTypes: []string{DigitalPinEventRisingEdge, DigitalPinEventRisingEdge},
		},
		"no_edge_rising": {
			eventOnEdge: digitalPinEventOnRisingEdge,
			simulateReadValues: []readValue{
				{value: 1},
				{value: 0},
				{value: 0},
			},
			wantEdgeTypes: nil,
		},
		"edge_both": {
			eventOnEdge: digitalPinEventOnBothEdges,
			simulateReadValues: []readValue{
				{value: 0},
				{value: 1},
				{value: 0},
				{value: 1},
				{value: 1},
			},
			wantEdgeTypes: []string{DigitalPinEventRisingEdge, DigitalPinEventFallingEdge, DigitalPinEventRisingEdge},
		},
		"no_edges_low": {
			eventOnEdge: digitalPinEventOnBothEdges,
			simulateReadValues: []readValue{
				{value: 0},
				{value: 0},
				{value: 0},
			},
			wantEdgeTypes: nil,
		},
		"no_edges_high": {
			eventOnEdge: digitalPinEventOnBothEdges,
			simulateReadValues: []readValue{
				{value: 1},
				{value: 1},
				{value: 1},
			},
			wantEdgeTypes: nil,
		},
		"read_error_keep_state": {
			eventOnEdge: digitalPinEventOnBothEdges,
			simulateReadValues: []readValue{
				{value: 0},
				{value: 1, err: "read error suppress rising and falling edge"},
				{value: 0},
				{value: 1},
				{value: 1},
			},
			wantEdgeTypes: []string{DigitalPinEventRisingEdge},
		},
		"error_no_eventhandler": {
			simulateNoEventHandler: true,
			wantErr:                "event handler is mandatory",
		},
		"error_no_quitchannel": {
			simulateNoQuitChan: true,
			wantErr:            "quit channel is mandatory",
		},
		"error_unsupported_edgetype_none": {
			eventOnEdge: digitalPinEventNone,
			wantErr:     "unsupported edge type 0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			pinLabel := "test_pin"
			pollInterval := time.Microsecond // zero is possible, just to show usage
			// arrange event handler
			var edgeTypes []string
			var eventHandler func(int, time.Duration, string, uint32, uint32)
			if !tc.simulateNoEventHandler {
				eventHandler = func(offset int, t time.Duration, et string, sn uint32, lsn uint32) {
					edgeTypes = append(edgeTypes, et)
				}
			}
			// arrange quit channel
			var quitChan chan struct{}
			if !tc.simulateNoQuitChan {
				quitChan = make(chan struct{})
			}
			defer func() {
				if quitChan != nil {
					close(quitChan)
				}
			}()
			// arrange reads
			numCallsRead := 0
			wg := sync.WaitGroup{}
			if tc.simulateReadValues != nil {
				wg.Add(1)
			}
			readFunc := func() (int, error) {
				numCallsRead++
				readVal := tc.simulateReadValues[numCallsRead-1]
				var err error
				if readVal.err != "" {
					err = errors.New(readVal.err)
				}
				if numCallsRead >= len(tc.simulateReadValues) {
					close(quitChan) // ensure no further read call
					quitChan = nil  // lets skip defer routine
					wg.Done()       // release assertions
				}

				return readVal.value, err
			}
			// act
			err := startEdgePolling(pinLabel, readFunc, pollInterval, tc.eventOnEdge, eventHandler, quitChan)
			wg.Wait()
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Len(t, tc.simulateReadValues, numCallsRead)
			assert.Equal(t, tc.wantEdgeTypes, edgeTypes)
		})
	}
}
