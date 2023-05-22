package joule

var sysfsPinMap = map[string]sysfsPin{
	// GPIO22
	"J12_1": {
		pin:    451,
		pwmPin: -1,
	},
	// SPP1RX
	"J12_2": {
		pin:    421,
		pwmPin: -1,
	},
	// PMICRST
	"J12_3": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP1TX
	"J12_4": {
		pin:    422,
		pwmPin: -1,
	},
	// 19.2mhz
	"J12_5": {
		pin:    356,
		pwmPin: -1,
	},
	// SPP1FS0
	"J12_6": {
		pin:    417,
		pwmPin: -1,
	},
	// UART0TX
	"J12_7": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP1FS2
	"J12_8": {
		pin:    419,
		pwmPin: -1,
	},
	// PWRGD
	"J12_9": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP1CLK
	"J12_10": {
		pin:    416,
		pwmPin: -1,
	},
	// I2C0SDA
	"J12_11": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2S1SDI
	"J12_12": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2C0SCL
	"J12_13": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2S1SDO
	"J12_14": {
		pin:    -1,
		pwmPin: -1,
	},
	// II0SDA
	"J12_15": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2S1WS
	"J12_16": {
		pin:    380,
		pwmPin: -1,
	},
	// IIC0SCL
	"J12_17": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2S1CLK
	"J12_18": {
		pin:    379,
		pwmPin: -1,
	},
	// IIC1SDA
	"J12_19": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2S1MCL
	"J12_20": {
		pin:    378,
		pwmPin: -1,
	},
	// IIC1SCL
	"J12_21": {
		pin:    -1,
		pwmPin: -1,
	},
	// UART1TX
	"J12_22": {
		pin:    -1,
		pwmPin: -1,
	},
	// ISH_IO6
	"J12_23": {
		pin:    343,
		pwmPin: -1,
	},
	// UART1RX
	"J12_24": {
		pin:    -1,
		pwmPin: -1,
	},
	// ISH_IO5
	"J12_25": {
		pin:    342,
		pwmPin: -1,
	},
	// PWM0
	"J12_26": {
		pin:    463,
		pwmPin: 0,
	},
	// ISH_IO4
	"J12_27": {
		pin:    341,
		pwmPin: -1,
	},
	// PWM1
	"J12_28": {
		pin:    464,
		pwmPin: 1,
	},
	// ISH_IO3
	"J12_29": {
		pin:    340,
		pwmPin: -1,
	},
	// PWM2
	"J12_30": {
		pin:    465,
		pwmPin: 2,
	},
	// ISH_IO2
	"J12_31": {
		pin:    339,
		pwmPin: -1,
	},
	// PWM3
	"J12_32": {
		pin:    466,
		pwmPin: 3,
	},
	// ISH_IO1
	"J12_33": {
		pin:    338,
		pwmPin: -1,
	},
	// 1.8V
	"J12_34": {
		pin:    -1,
		pwmPin: -1,
	},
	// ISH_IO0
	"J12_35": {
		pin:    337,
		pwmPin: -1,
	},
	// GND
	"J12_36": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J12_37": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J12_38": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J12_39": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J12_40": {
		pin:    -1,
		pwmPin: -1,
	},

	// Second header
	// GND
	"J13_1": {
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"J13_2": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J13_3": {
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"J13_4": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J13_5": {
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"J13_6": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J13_7": {
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"J13_8": {
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"J13_9": {
		pin:    -1,
		pwmPin: -1,
	},
	// 1.8V
	"J13_10": {
		pin:    -1,
		pwmPin: -1,
	},
	// GPIO
	"J13_11": {
		pin:    456,
		pwmPin: -1,
	},
	// 1.8V
	"J13_12": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"J13_13": {
		pin:    270,
		pwmPin: -1,
	},
	// GND
	"J13_14": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"J13_15": {
		pin:    271,
		pwmPin: -1,
	},
	// CAMERA
	"J13_16": {
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"J13_17": {
		pin:    272,
		pwmPin: -1,
	},
	// CAMERA
	"J13_18": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS0
	"J13_19": {
		pin:    411,
		pwmPin: -1,
	},
	// CAMERA
	"J13_20": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS1
	"J13_21": {
		pin:    412,
		pwmPin: -1,
	},
	// SPI_DAT
	"J13_22": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS2
	"J13_23": {
		pin:    413,
		pwmPin: -1,
	},
	// SPICLKB
	"J13_24": {
		pin:    384,
		pwmPin: -1,
	},
	// SPP0CLK
	"J13_25": {
		pin:    410,
		pwmPin: -1,
	},
	// SPICLKA
	"J13_26": {
		pin:    383,
		pwmPin: -1,
	},
	// SPP0TX
	"J13_27": {
		pin:    414,
		pwmPin: -1,
	},
	// UART0RX
	"J13_28": {
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0RX
	"J13_29": {
		pin:    415,
		pwmPin: -1,
	},
	// UART0RT
	"J13_30": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2C1SDA
	"J13_31": {
		pin:    -1,
		pwmPin: -1,
	},
	// UART0CT
	"J13_32": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2C1SCL
	"J13_33": {
		pin:    -1,
		pwmPin: -1,
	},
	// IURT0TX
	"J13_34": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2C2SDA
	"J13_35": {
		pin:    -1,
		pwmPin: -1,
	},
	// IURT0RX
	"J13_36": {
		pin:    -1,
		pwmPin: -1,
	},
	// I2C2SCL
	"J13_37": {
		pin:    -1,
		pwmPin: -1,
	},
	// IURT0RT
	"J13_38": {
		pin:    -1,
		pwmPin: -1,
	},
	// RTC_CLK
	"J13_39": {
		pin:    367,
		pwmPin: -1,
	},
	// IURT0CT
	"J13_40": {
		pin:    -1,
		pwmPin: -1,
	},

	// Built-in LEDs
	// LED100
	"GP100": {
		pin:    337,
		pwmPin: -1,
	},
	// LED101
	"GP101": {
		pin:    338,
		pwmPin: -1,
	},
	// LED102
	"GP102": {
		pin:    339,
		pwmPin: -1,
	},
	// LED103
	"GP103": {
		pin:    340,
		pwmPin: -1,
	},
	// LEDWIFI
	"GP104": {
		pin:    438,
		pwmPin: -1,
	},
	// LEDBT
	"GP105": {
		pin:    439,
		pwmPin: -1,
	},
}
