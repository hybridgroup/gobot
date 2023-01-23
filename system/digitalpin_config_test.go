package system

import (
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.label, label)
	gobottest.Assert(t, d.direction, IN)
	gobottest.Assert(t, d.outInitialState, 0)
}

func Test_newDigitalPinConfigWithOption(t *testing.T) {
	// arrange
	const label = "gobotio18"
	// act
	d := newDigitalPinConfig("not used", WithPinLabel(label))
	// assert
	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.label, label)
}

func TestWithPinLabel(t *testing.T) {
	const (
		oldLabel = "old label"
		newLabel = "my optional label"
	)
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.label, tc.setLabel)
		})
	}
}

func TestWithPinDirectionOutput(t *testing.T) {
	const (
		// values other than 0, 1 are normally not useful, just to test
		oldVal = 3
		newVal = 5
	)
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.direction, "out")
			gobottest.Assert(t, dpc.outInitialState, tc.wantVal)
		})
	}
}

func TestWithPinDirectionInput(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.direction, "in")
			gobottest.Assert(t, dpc.outInitialState, initValOut)
		})
	}
}

func TestWithPinActiveLow(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.activeLow, true)
		})
	}
}

func TestWithPinPullDown(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.bias, digitalPinBiasPullDown)
		})
	}
}

func TestWithPinPullUp(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.bias, digitalPinBiasPullUp)
		})
	}
}

func TestWithPinOpenDrain(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.drive, digitalPinDriveOpenDrain)
		})
	}
}

func TestWithPinOpenSource(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.drive, digitalPinDriveOpenSource)
		})
	}
}

func TestWithPinDebounce(t *testing.T) {
	const (
		oldVal = time.Duration(10)
		newVal = time.Duration(14)
	)
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.debouncePeriod, newVal)
		})
	}
}
