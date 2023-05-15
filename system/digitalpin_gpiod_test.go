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
	// This is a general test, that options are applied by using "newDigitalPinGpiod" with the WithPinLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my own label"
	// act
	dp := newDigitalPinGpiod("", 9, WithPinLabel(label))
	// assert
	gobottest.Assert(t, dp.label, label)
}

func TestApplyOptions(t *testing.T) {
	var tests = map[string]struct {
		changed          []bool
		simErr           error
		wantReconfigured int
		wantErr          error
	}{
		"both_changed": {
			changed:          []bool{true, true},
			wantReconfigured: 1,
		},
		"first_changed": {
			changed:          []bool{true, false},
			wantReconfigured: 1,
		},
		"second_changed": {
			changed:          []bool{false, true},
			wantReconfigured: 1,
		},
		"none_changed": {
			changed:          []bool{false, false},
			simErr:           fmt.Errorf("error not raised"),
			wantReconfigured: 0,
		},
		"error_on_change": {
			changed:          []bool{false, true},
			simErr:           fmt.Errorf("error raised"),
			wantReconfigured: 1,
			wantErr:          fmt.Errorf("error raised"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// currently the gpiod.Chip has no interface for RequestLine(),
			// so we can only test without trigger of real reconfigure
			// arrange
			orgReconf := digitalPinGpiodReconfigure
			defer func() { digitalPinGpiodReconfigure = orgReconf }()

			inputForced := true
			reconfigured := 0
			digitalPinGpiodReconfigure = func(d *digitalPinGpiod, forceInput bool) error {
				inputForced = forceInput
				reconfigured++
				return tc.simErr
			}
			d := &digitalPinGpiod{digitalPinConfig: &digitalPinConfig{direction: "in"}}
			optionFunction1 := func(gobot.DigitalPinOptioner) bool {
				d.digitalPinConfig.direction = "test"
				return tc.changed[0]
			}
			optionFunction2 := func(gobot.DigitalPinOptioner) bool {
				d.digitalPinConfig.drive = 15
				return tc.changed[1]
			}
			// act
			err := d.ApplyOptions(optionFunction1, optionFunction2)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, d.digitalPinConfig.direction, "test")
			gobottest.Assert(t, d.digitalPinConfig.drive, 15)
			gobottest.Assert(t, reconfigured, tc.wantReconfigured)
			if reconfigured > 0 {
				gobottest.Assert(t, inputForced, false)
			}
		})
	}
}

func TestExport(t *testing.T) {
	var tests = map[string]struct {
		simErr           error
		wantReconfigured int
		wantErr          error
	}{
		"no_err": {
			wantReconfigured: 1,
		},
		"error": {
			wantReconfigured: 1,
			simErr:           fmt.Errorf("reconfigure error"),
			wantErr:          fmt.Errorf("gpiod.Export(): reconfigure error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// currently the gpiod.Chip has no interface for RequestLine(),
			// so we can only test without trigger of real reconfigure
			// arrange
			orgReconf := digitalPinGpiodReconfigure
			defer func() { digitalPinGpiodReconfigure = orgReconf }()

			inputForced := true
			reconfigured := 0
			digitalPinGpiodReconfigure = func(d *digitalPinGpiod, forceInput bool) error {
				inputForced = forceInput
				reconfigured++
				return tc.simErr
			}
			d := &digitalPinGpiod{}
			// act
			err := d.Export()
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, inputForced, false)
			gobottest.Assert(t, reconfigured, tc.wantReconfigured)
		})
	}
}

func TestUnexport(t *testing.T) {
	var tests = map[string]struct {
		simNoLine        bool
		simReconfErr     error
		simCloseErr      error
		wantReconfigured int
		wantErr          error
	}{
		"no_line_no_err": {
			simNoLine:        true,
			wantReconfigured: 0,
		},
		"no_line_with_err": {
			simNoLine:        true,
			simReconfErr:     fmt.Errorf("reconfigure error"),
			wantReconfigured: 0,
		},
		"no_err": {
			wantReconfigured: 1,
		},
		"error_reconfigure": {
			wantReconfigured: 1,
			simReconfErr:     fmt.Errorf("reconfigure error"),
			wantErr:          fmt.Errorf("reconfigure error"),
		},
		"error_close": {
			wantReconfigured: 1,
			simCloseErr:      fmt.Errorf("close error"),
			wantErr:          fmt.Errorf("gpiod.Unexport()-line.Close(): close error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// currently the gpiod.Chip has no interface for RequestLine(),
			// so we can only test without trigger of real reconfigure
			// arrange
			orgReconf := digitalPinGpiodReconfigure
			defer func() { digitalPinGpiodReconfigure = orgReconf }()

			inputForced := false
			reconfigured := 0
			digitalPinGpiodReconfigure = func(d *digitalPinGpiod, forceInput bool) error {
				inputForced = forceInput
				reconfigured++
				return tc.simReconfErr
			}
			dp := newDigitalPinGpiod("", 4)
			if !tc.simNoLine {
				dp.line = &lineMock{simCloseErr: tc.simCloseErr}
			}
			// act
			err := dp.Unexport()
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, reconfigured, tc.wantReconfigured)
			if reconfigured > 0 {
				gobottest.Assert(t, inputForced, true)
			}
		})
	}
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
			lm := &lineMock{lastVal: 10, simSetErr: tc.simErr}
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
			lm := &lineMock{lastVal: tc.simVal, simValueErr: tc.simErr}
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
	lastVal     int
	simSetErr   error
	simValueErr error
	simCloseErr error
}

func (lm *lineMock) SetValue(value int) error { lm.lastVal = value; return lm.simSetErr }
func (lm *lineMock) Value() (int, error)      { return lm.lastVal, lm.simValueErr }
func (lm *lineMock) Close() error             { return lm.simCloseErr }
