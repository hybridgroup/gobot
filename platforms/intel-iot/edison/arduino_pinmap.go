package edison

import "gobot.io/x/gobot/sysfs"

var arduinoPinMap = map[string]sysfsPin{
	"0": {
		pin:          130,
		resistor:     216,
		levelShifter: 248,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"1": {
		pin:          131,
		resistor:     217,
		levelShifter: 249,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"2": {
		pin:          128,
		resistor:     218,
		levelShifter: 250,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"3": {
		pin:          12,
		resistor:     219,
		levelShifter: 251,
		pwmPin:       0,
		mux:          []mux{},
	},

	"4": {
		pin:          129,
		resistor:     220,
		levelShifter: 252,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"5": {
		pin:          13,
		resistor:     221,
		levelShifter: 253,
		pwmPin:       1,
		mux:          []mux{},
	},
	"6": {
		pin:          182,
		resistor:     222,
		levelShifter: 254,
		pwmPin:       2,
		mux:          []mux{},
	},
	"7": {
		pin:          48,
		resistor:     223,
		levelShifter: 255,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"8": {
		pin:          49,
		resistor:     224,
		levelShifter: 256,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"9": {
		pin:          183,
		resistor:     225,
		levelShifter: 257,
		pwmPin:       3,
		mux:          []mux{},
	},
	"10": {
		pin:          41,
		resistor:     226,
		levelShifter: 258,
		pwmPin:       4,
		mux: []mux{
			{263, sysfs.HIGH},
			{240, sysfs.LOW},
		},
	},
	"11": {
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			{262, sysfs.HIGH},
			{241, sysfs.LOW},
		},
	},
	"12": {
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			{242, sysfs.LOW},
		},
	},
	"13": {
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			{243, sysfs.LOW},
		},
	},
}
