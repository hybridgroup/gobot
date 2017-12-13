package up2

var fixedPins = map[string]sysfsPin{
	"7": {
		pin:    462, // GPIO4
		pwmPin: -1,
	},
	"13": {
		pin:    432, // GPIO27
		pwmPin: -1,
	},
	"15": {
		pin:    431, // GPIO22
		pwmPin: -1,
	},
	"16": {
		pin:    471, // PWM3
		pwmPin: 3,
	},
	"18": {
		pin:    405, // GPIO24
		pwmPin: -1,
	},
	"22": {
		pin:    402, // GPIO25
		pwmPin: -1,
	},
	"29": {
		pin:    430, // GPIO5
		pwmPin: -1,
	},
	"31": {
		pin:    404, // GPIO6
		pwmPin: -1,
	},
	"32": {
		pin:    468, // PWM0
		pwmPin: 0,
	},
	"33": {
		pin:    469, // PWM1
		pwmPin: 1,
	},
	"37": {
		pin:    403, // GPIO26
		pwmPin: -1,
	},
}
