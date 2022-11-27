package system

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.DigitalPinner = (*digitalPinGpiod)(nil)
var _ gobot.DigitalPinValuer = (*digitalPinGpiod)(nil)
var _ gobot.DigitalPinOptioner = (*digitalPinGpiod)(nil)
var _ gobot.DigitalPinOptionApplier = (*digitalPinGpiod)(nil)

func Test_newDigitalPinGpiod(t *testing.T) {
	// arrange
	const (
		chip  = "gpiochip0"
		pin   = 17
		label = "gobotio17"
	)
	// act
	d := newDigitalPinGpiod(chip, pin)
	// assert
	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.chipName, chip)
	gobottest.Assert(t, d.pin, pin)
	gobottest.Assert(t, d.label, label)
	gobottest.Assert(t, d.direction, IN)
	gobottest.Assert(t, d.outInitialState, 0)
}

func Test_newDigitalPinGpiodWithOptions(t *testing.T) {
	// This is a general test, that options are applied by using "newDigitalPinGpiod" with the WithLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my own label"
	// act
	dp := newDigitalPinGpiod("", 9, WithLabel(label))
	// assert
	gobottest.Assert(t, dp.label, label)
}

func TestApplyOptions(t *testing.T) {
	// currently the gpiod.Chip has no interface for RequestLine(),
	// so we can only test without trigger of reconfigure
	// arrange
	d := &digitalPinGpiod{digitalPinConfig: &digitalPinConfig{direction: "in"}}
	// act
	d.ApplyOptions(WithDirectionInput())
	// assert
	gobottest.Assert(t, d.digitalPinConfig.direction, "in")
}

func TestUnexport(t *testing.T) {
	// currently the gpiod.Chip has no interface for RequestLine(),
	// so we can only test without trigger of reconfigure
	// arrange
	dp := newDigitalPinGpiod("", 4)
	dp.line = nil // ensures no reconfigure
	// act
	err := dp.Unexport()
	// assert
	gobottest.Assert(t, err, nil)
}

func TestWrite(t *testing.T) {
	var tests = map[string]struct {
		val     int
		simErr  error
		want    int
		wantErr []string
	}{
		"write_zero": {
			val:  0,
			want: 0,
		},
		"write_one": {
			val:  1,
			want: 1,
		},
		"write_minus_one": {
			val:  -1,
			want: 0,
		},
		"write_two": {
			val:  2,
			want: 1,
		},
		"write_with_err": {
			simErr:  fmt.Errorf("a write err"),
			wantErr: []string{"a write err", "gpiod.Write"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dp := newDigitalPinGpiod("", 4)
			lm := &lineMock{lastVal: 10, simErr: tc.simErr}
			dp.line = lm
			// act
			err := dp.Write(tc.val)
			// assert
			if tc.wantErr != nil {
				for _, want := range tc.wantErr {
					gobottest.Assert(t, strings.Contains(err.Error(), want), true)
				}
			} else {
				gobottest.Assert(t, err, nil)
			}
			gobottest.Assert(t, lm.lastVal, tc.want)
		})
	}
}

func TestRead(t *testing.T) {
	var tests = map[string]struct {
		simVal  int
		simErr  error
		wantErr []string
	}{
		"read_ok": {
			simVal: 3,
		},
		"write_with_err": {
			simErr:  fmt.Errorf("a read err"),
			wantErr: []string{"a read err", "gpiod.Read"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dp := newDigitalPinGpiod("", 4)
			lm := &lineMock{lastVal: tc.simVal, simErr: tc.simErr}
			dp.line = lm
			// act
			got, err := dp.Read()
			// assert
			if tc.wantErr != nil {
				for _, want := range tc.wantErr {
					gobottest.Assert(t, strings.Contains(err.Error(), want), true)
				}
			} else {
				gobottest.Assert(t, err, nil)
			}
			gobottest.Assert(t, tc.simVal, got)
		})
	}
}

type lineMock struct {
	lastVal int
	simErr  error
}

func (lm *lineMock) SetValue(value int) error { lm.lastVal = value; return lm.simErr }
func (lm *lineMock) Value() (int, error)      { return lm.lastVal, lm.simErr }
func (*lineMock) Close() error                { return nil }
