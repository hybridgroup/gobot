package joule

var sysfsPinMap = map[string]sysfsPin{
	// disabled
	"0": {
		pin:    -1,
		pwmPin: -1,
	},
	"1": {
		pin:    446,
		pwmPin: -1,
	},
	"2": {
		pin:    421,
		pwmPin: -1,
	},
	// disabled
	"3": {
		pin:    -1,
		pwmPin: -1,
	},
	"4": {
		pin:    422,
		pwmPin: -1,
	},
	"5": {
		pin:    356,
		pwmPin: -1,
	},
	"6": {
		pin:    417,
		pwmPin: -1,
	},
	// UART
	"7": {
		pin:    -1,
		pwmPin: -1,
	},
	"8": {
		pin:    419,
		pwmPin: -1,
	},
	// disabled
	"9": {
		pin:    -1,
		pwmPin: -1,
	},
	"10": {
		pin:    416,
		pwmPin: -1,
	},
	"11": {
		pin:    381,
		pwmPin: -1,
	},
	"13": {
		pin:    382,
		pwmPin: -1,
	},
	"15": {
		pin:    380,
		pwmPin: -1,
	},
	"17": {
		pin:    379,
		pwmPin: -1,
	},
	"19": {
		pin:    378,
		pwmPin: -1,
	},
	// UART
	"21": {
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"22": {
		pin:    -1,
		pwmPin: -1,
	},
	// UART
	"23": {
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"24": {
		pin:    -1,
		pwmPin: -1,
	},
	"25": {
		pin:    463,
		pwmPin: 0,
	},
	// low voltage should not use
	"26": {
		pin:    -1,
		pwmPin: -1,
	},
	"27": {
		pin:    464,
		pwmPin: 1,
	},
	// disabled
	"28": {
		pin:    -1,
		pwmPin: -1,
	},
	"29": {
		pin:    465,
		pwmPin: 2,
	},
	// disabled?
	"30": {
		pin:    -1,
		pwmPin: -1,
	},
	"31": {
		pin:    466,
		pwmPin: 3,
	},
	// disabled?
	"32": {
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"33": {
		pin:    -1,
		pwmPin: -1,
	},
	"34": {
		pin:    393,
		pwmPin: -1,
	},
	// GND
	"35": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"36": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"37": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"38": {
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"39": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"40": {
		pin:    -1,
		pwmPin: -1,
	},

	// Second header
	// GND
	"41": {
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"42": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"43": {
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"44": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"45": {
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"46": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"47": {
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"48": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"49": {
		pin:    -1,
		pwmPin: -1,
	},
	// 1.8V
	"50": {
		pin:    -1,
		pwmPin: -1,
	},
	// GPIO
	"51": {
		pin:    456,
		pwmPin: -1,
	},
	// 1.8V
	"52": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"53": {
		pin:    270,
		pwmPin: -1,
	},
	// GND
	"54": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"55": {
		pin:    271,
		pwmPin: -1,
	},
	// CAMERA
	"56": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"57": {
		pin:    272,
		pwmPin: -1,
	},
	// CAMERA
	"58": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS0
	"59": {
		pin:    411,
		pwmPin: -1,
	},
	// CAMERA
	"60": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS1
	"61": {
		pin:    412,
		pwmPin: -1,
	},
	// SPI_DAT
	"62": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS2
	"63": {
		pin:    411,
		pwmPin: -1,
	},
	// SPICLKB
	"64": {
		pin:    384,
		pwmPin: -1,
	},
	// SPP0FS3
	"65": {
		pin:    410,
		pwmPin: -1,
	},
	// SPICLKA
	"66": {
		pin:    383,
		pwmPin: -1,
	},
	// SPP0TX
	"67": {
		pin:    414,
		pwmPin: -1,
	},
	// UART0RX
	"68": {
		pin:    467,
		pwmPin: -1,
	},
	// SPP0RX
	"69": {
		pin:    415,
		pwmPin: -1,
	},
	// UART0RT
	"70": {
		pin:    469,
		pwmPin: -1,
	},
	// I2C1SDA
	"71": {
		pin:    317,
		pwmPin: -1,
	},
	// UART0CT
	"72": {
		pin:    412,
		pwmPin: -1,
	},
	// I2C1SCL
	"73": {
		pin:    418,
		pwmPin: -1,
	},
	// UART1TX
	"74": {
		pin:    484,
		pwmPin: -1,
	},
	// I2C2SDA
	"75": {
		pin:    319,
		pwmPin: -1,
	},
	// UART1RX
	"76": {
		pin:    483,
		pwmPin: -1,
	},
	// I2C2SCL
	"77": {
		pin:    320,
		pwmPin: -1,
	},
	// UART1RT
	"78": {
		pin:    485,
		pwmPin: -1,
	},
	// RTC_CLK
	"79": {
		pin:    367,
		pwmPin: -1,
	},
	// UART1CT
	"80": {
		pin:    486,
		pwmPin: -1,
	},

	// Built-in LEDs
	// LED100
	"100": {
		pin:    337,
		pwmPin: -1,
	},
	// LED101
	"101": {
		pin:    338,
		pwmPin: -1,
	},
	// LED102
	"102": {
		pin:    339,
		pwmPin: -1,
	},
	// LED103
	"103": {
		pin:    340,
		pwmPin: -1,
	},
	// LEDWIFI
	"104": {
		pin:    438,
		pwmPin: -1,
	},
	// LEDBT
	"105": {
		pin:    439,
		pwmPin: -1,
	},
}
