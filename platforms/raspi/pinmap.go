package raspi

import "gobot.io/x/gobot/v2/platforms/adaptors"

var pins = map[string]map[string]int{
	"3": {
		"1": 0,
		"2": 2,
		"3": 2,
	},
	"5": {
		"1": 1,
		"2": 3,
		"3": 3,
	},
	"7": {
		"*": 4,
	},
	"8": {
		"*": 14,
	},
	"10": {
		"*": 15,
	},
	"11": {
		"*": 17,
	},
	"12": {
		"*": 18,
	},
	"13": {
		"1": 21,
		"2": 27,
		"3": 27,
	},
	"15": {
		"*": 22,
	},
	"16": {
		"*": 23,
	},
	"18": {
		"*": 24,
	},
	"19": {
		"*": 10,
	},
	"21": {
		"*": 9,
	},
	"22": {
		"*": 25,
	},
	"23": {
		"*": 11,
	},
	"24": {
		"*": 8,
	},
	"26": {
		"*": 7,
	},
	"29": {
		"3": 5,
	},
	"31": {
		"3": 6,
	},
	"32": {
		"3": 12,
	},
	"33": {
		"3": 13,
	},
	"35": {
		"3": 19,
	},
	"36": {
		"3": 16,
	},
	"37": {
		"3": 26,
	},
	"38": {
		"3": 20,
	},
	"40": {
		"3": 21,
	},
	"pwm0": { // pin 12 (GPIO18) and pin 32 (GPIO12) can be configured for "pwm0"
		"*": 0,
	},
	"pwm1": { // pin 33 (GPIO13) and pin 35 (GPIO19) can be configured for "pwm1"
		"3": 1,
	},
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	"thermal_zone0": {Path: "/sys/class/thermal/thermal_zone0/temp", W: false, ReadBufLen: 7},
}
