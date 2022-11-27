package system

import (
	"testing"

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

func TestWithLabel(t *testing.T) {
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
			got := WithLabel(tc.setLabel)(dpc)
			// assert
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.label, tc.setLabel)
		})
	}
}

func TestWithDirectionOutput(t *testing.T) {
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
			got := WithDirectionOutput(newVal)(dpc)
			// assert
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.direction, "out")
			gobottest.Assert(t, dpc.outInitialState, tc.wantVal)
		})
	}
}

func TestWithDirectionInput(t *testing.T) {
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
			got := WithDirectionInput()(dpc)
			// assert
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, dpc.direction, "in")
			gobottest.Assert(t, dpc.outInitialState, initValOut)
		})
	}
}
