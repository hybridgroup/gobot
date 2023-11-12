package system

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var _ gobot.DigitalPinOptioner = (*digitalPinConfig)(nil)

func Test_newDigitalPinConfig(t *testing.T) {
	// arrange
	const (
		label = "gobotio17"
	)
	// act
	d := newDigitalPinConfig(label)
	// assert
	assert.NotNil(t, d)
	assert.Equal(t, label, d.label)
	assert.Equal(t, IN, d.direction)
	assert.Equal(t, 0, d.outInitialState)
}

func Test_newDigitalPinConfigWithOption(t *testing.T) {
	// arrange
	const label = "gobotio18"
	// act
	d := newDigitalPinConfig("not used", WithPinLabel(label))
	// assert
	assert.NotNil(t, d)
	assert.Equal(t, label, d.label)
}

func TestWithPinLabel(t *testing.T) {
	const (
		oldLabel = "old label"
		newLabel = "my optional label"
	)
	tests := map[string]struct {
		setLabel string
		want     bool
	}{
		"no_change": {
			setLabel: oldLabel,
		},
		"change": {
			setLabel: newLabel,
			want:     true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{label: oldLabel}
			// act
			got := WithPinLabel(tc.setLabel)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.setLabel, dpc.label)
		})
	}
}

func TestWithPinDirectionOutput(t *testing.T) {
	const (
		// values other than 0, 1 are normally not useful, just to test
		oldVal = 3
		newVal = 5
	)
	tests := map[string]struct {
		oldDir  string
		want    bool
		wantVal int
	}{
		"no_change": {
			oldDir:  "out",
			wantVal: oldVal,
		},
		"change": {
			oldDir:  "in",
			want:    true,
			wantVal: newVal,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{direction: tc.oldDir, outInitialState: oldVal}
			// act
			got := WithPinDirectionOutput(newVal)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, "out", dpc.direction)
			assert.Equal(t, tc.wantVal, dpc.outInitialState)
		})
	}
}

func TestWithPinDirectionInput(t *testing.T) {
	tests := map[string]struct {
		oldDir string
		want   bool
	}{
		"no_change": {
			oldDir: "in",
		},
		"change": {
			oldDir: "out",
			want:   true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const initValOut = 2 // 2 is normally not useful, just to test that is not touched
			dpc := &digitalPinConfig{direction: tc.oldDir, outInitialState: initValOut}
			// act
			got := WithPinDirectionInput()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, "in", dpc.direction)
			assert.Equal(t, initValOut, dpc.outInitialState)
		})
	}
}

func TestWithPinActiveLow(t *testing.T) {
	tests := map[string]struct {
		oldActiveLow bool
		want         bool
	}{
		"no_change": {
			oldActiveLow: true,
		},
		"change": {
			oldActiveLow: false,
			want:         true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{activeLow: tc.oldActiveLow}
			// act
			got := WithPinActiveLow()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.True(t, dpc.activeLow)
		})
	}
}

func TestWithPinPullDown(t *testing.T) {
	tests := map[string]struct {
		oldBias int
		want    bool
		wantVal int
	}{
		"no_change": {
			oldBias: digitalPinBiasPullDown,
		},
		"change": {
			oldBias: digitalPinBiasPullUp,
			want:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{bias: tc.oldBias}
			// act
			got := WithPinPullDown()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, digitalPinBiasPullDown, dpc.bias)
		})
	}
}

func TestWithPinPullUp(t *testing.T) {
	tests := map[string]struct {
		oldBias int
		want    bool
		wantVal int
	}{
		"no_change": {
			oldBias: digitalPinBiasPullUp,
		},
		"change": {
			oldBias: digitalPinBiasPullDown,
			want:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{bias: tc.oldBias}
			// act
			got := WithPinPullUp()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, digitalPinBiasPullUp, dpc.bias)
		})
	}
}

func TestWithPinOpenDrain(t *testing.T) {
	tests := map[string]struct {
		oldDrive int
		want     bool
		wantVal  int
	}{
		"no_change": {
			oldDrive: digitalPinDriveOpenDrain,
		},
		"change_from_pushpull": {
			oldDrive: digitalPinDrivePushPull,
			want:     true,
		},
		"change_from_opensource": {
			oldDrive: digitalPinDriveOpenSource,
			want:     true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{drive: tc.oldDrive}
			// act
			got := WithPinOpenDrain()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, digitalPinDriveOpenDrain, dpc.drive)
		})
	}
}

func TestWithPinOpenSource(t *testing.T) {
	tests := map[string]struct {
		oldDrive int
		want     bool
		wantVal  int
	}{
		"no_change": {
			oldDrive: digitalPinDriveOpenSource,
		},
		"change_from_pushpull": {
			oldDrive: digitalPinDrivePushPull,
			want:     true,
		},
		"change_from_opendrain": {
			oldDrive: digitalPinDriveOpenDrain,
			want:     true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{drive: tc.oldDrive}
			// act
			got := WithPinOpenSource()(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, digitalPinDriveOpenSource, dpc.drive)
		})
	}
}

func TestWithPinDebounce(t *testing.T) {
	const (
		oldVal = time.Duration(10)
		newVal = time.Duration(14)
	)
	tests := map[string]struct {
		oldDebouncePeriod time.Duration
		want              bool
		wantVal           time.Duration
	}{
		"no_change": {
			oldDebouncePeriod: newVal,
		},
		"change": {
			oldDebouncePeriod: oldVal,
			want:              true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{debouncePeriod: tc.oldDebouncePeriod}
			// act
			got := WithPinDebounce(newVal)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, newVal, dpc.debouncePeriod)
		})
	}
}

func TestWithPinEventOnFallingEdge(t *testing.T) {
	const (
		oldVal = digitalPinEventNone
		newVal = digitalPinEventOnFallingEdge
	)
	tests := map[string]struct {
		oldEdge int
		want    bool
		wantVal int
	}{
		"no_change": {
			oldEdge: newVal,
			want:    false,
		},
		"change": {
			oldEdge: oldVal,
			want:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{edge: tc.oldEdge}
			handler := func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32, lseqno uint32) {}
			// act
			got := WithPinEventOnFallingEdge(handler)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, newVal, dpc.edge)
			if tc.want {
				assert.NotNil(t, dpc.edgeEventHandler)
			} else {
				assert.Nil(t, dpc.edgeEventHandler)
			}
		})
	}
}

func TestWithPinEventOnRisingEdge(t *testing.T) {
	const (
		oldVal = digitalPinEventNone
		newVal = digitalPinEventOnRisingEdge
	)
	tests := map[string]struct {
		oldEdge int
		want    bool
		wantVal int
	}{
		"no_change": {
			oldEdge: newVal,
			want:    false,
		},
		"change": {
			oldEdge: oldVal,
			want:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{edge: tc.oldEdge}
			handler := func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32, lseqno uint32) {}
			// act
			got := WithPinEventOnRisingEdge(handler)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, newVal, dpc.edge)
			if tc.want {
				assert.NotNil(t, dpc.edgeEventHandler)
			} else {
				assert.Nil(t, dpc.edgeEventHandler)
			}
		})
	}
}

func TestWithPinEventOnBothEdges(t *testing.T) {
	const (
		oldVal = digitalPinEventNone
		newVal = digitalPinEventOnBothEdges
	)
	tests := map[string]struct {
		oldEdge int
		want    bool
		wantVal int
	}{
		"no_change": {
			oldEdge: newVal,
			want:    false,
		},
		"change": {
			oldEdge: oldVal,
			want:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{edge: tc.oldEdge}
			handler := func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32, lseqno uint32) {}
			// act
			got := WithPinEventOnBothEdges(handler)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, newVal, dpc.edge)
			if tc.want {
				assert.NotNil(t, dpc.edgeEventHandler)
			} else {
				assert.Nil(t, dpc.edgeEventHandler)
			}
		})
	}
}

func TestWithPinPollForEdgeDetection(t *testing.T) {
	const (
		oldVal = time.Duration(1)
		newVal = time.Duration(3)
	)
	tests := map[string]struct {
		oldPollInterval time.Duration
		want            bool
		wantVal         time.Duration
	}{
		"no_change": {
			oldPollInterval: newVal,
		},
		"change": {
			oldPollInterval: oldVal,
			want:            true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dpc := &digitalPinConfig{pollInterval: tc.oldPollInterval}
			stopChan := make(chan struct{})
			defer close(stopChan)
			// act
			got := WithPinPollForEdgeDetection(newVal, stopChan)(dpc)
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, newVal, dpc.pollInterval)
		})
	}
}
