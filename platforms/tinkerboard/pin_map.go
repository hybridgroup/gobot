package tinkerboard

var fixedPins = map[string]sysfsPin{
	"7": {
		pin:    17, // GPIO0_C1
		pwmPin: -1,
	},
	"10": {
		pin:    160, // GPIO5_B0
		pwmPin: -1,
	},
	"8": {
		pin:    161, // GPIO5_B1
		pwmPin: -1,
	},
	"16": {
		pin:    162, // GPIO5_B2
		pwmPin: -1,
	},
	"18": {
		pin:    163, // GPIO5_B3
		pwmPin: -1,
	},
	"11": {
		pin:    164, // GPIO5_B4
		pwmPin: -1,
	},
	"29": {
		pin:    165, // GPIO5_B5
		pwmPin: -1,
	},
	"13": {
		pin:    166, // GPIO5_B6
		pwmPin: -1,
	},
	"15": {
		pin:    167, // GPIO5_B7
		pwmPin: -1,
	},
	"31": {
		pin:    168, // GPIO5_C0
		pwmPin: -1,
	},
	"22": {
		pin:    171, // GPIO5_C3
		pwmPin: -1,
	},
	"12": {
		pin:    184, // GPIO5_A0
		pwmPin: -1,
	},
	"35": {
		pin:    185, // GPIO5_A1
		pwmPin: -1,
	},
	"38": {
		pin:    187, // GPIO5_A3
		pwmPin: -1,
	},
	"40": {
		pin:    188, // GPIO5_A4
		pwmPin: -1,
	},
	"36": {
		pin:    223, // GPIO5_A7
		pwmPin: -1,
	},
	"37": {
		pin:    224, // GPIO5_B0
		pwmPin: -1,
	},
	"27": {
		pin:    233, // GPIO5_C1
		pwmPin: -1,
	},
	"28": {
		pin:    234, // GPIO5_C2
		pwmPin: -1,
	},
	"33": {
		pin:    238, // GPIO5_C6
		pwmPin: 0,
	},
	"32": {
		pin:    239, // GPIO5_C7
		pwmPin: 1,
	},
	"26": {
		pin:    251, // GPIO5_A3
		pwmPin: -1,
	},
	"3": {
		pin:    252, // GPIO5_A4
		pwmPin: -1,
	},
	"5": {
		pin:    253, // GPIO5_A3
		pwmPin: -1,
	},
	"23": {
		pin:    254, // GPIO5_A6
		pwmPin: -1,
	},
	"24": {
		pin:    255, // GPIO5_A7
		pwmPin: -1,
	},
	"21": {
		pin:    256, // GPIO5_B0
		pwmPin: -1,
	},
	"19": {
		pin:    257, // GPIO5_B1
		pwmPin: -1,
	},
}
