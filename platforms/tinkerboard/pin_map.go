package tinkerboard

var fixedPins = map[string]sysfsPin{
	"7": {
		pin:    17, // GPIO0_C1_CLKOUT
		pwmPin: -1,
	},
	"10": {
		pin:    160, // GPIO5_B0_UART1RX
		pwmPin: -1,
	},
	"8": {
		pin:    161, // GPIO5_B1_UART1TX
		pwmPin: -1,
	},
	"16": {
		pin:    162, // GPIO5_B2_UART1CTSN
		pwmPin: -1,
	},
	"18": {
		pin:    163, // GPIO5_B3_UART1RTSN
		pwmPin: -1,
	},
	"11": {
		pin:    164, // GPIO5_B4_SPI0CLK_UART4CTSN
		pwmPin: -1,
	},
	"29": {
		pin:    165, // GPIO5_B5_SPI0CSN_UART4RTSN
		pwmPin: -1,
	},
	"13": {
		pin:    166, // GPIO5_B6_SPI0_TXD_UART4TX
		pwmPin: -1,
	},
	"15": {
		pin:    167, // GPIO5_B7_SPI0_RXD_UART4RX
		pwmPin: -1,
	},
	"31": {
		pin:    168, // GPIO5_C0_SPI0CSN1
		pwmPin: -1,
	},
	"22": {
		pin:    171, // GPIO5_C3
		pwmPin: -1,
	},
	"12": {
		pin:    184, // GPIO6_A0_PCM/I2S_CLK
		pwmPin: -1,
	},
	"35": {
		pin:    185, // GPIO6_A1_PCM/I2S_FS
		pwmPin: -1,
	},
	"38": {
		pin:    187, // GPIO6_A3_PCM/I2S_SDI
		pwmPin: -1,
	},
	"40": {
		pin:    188, // GPIO6_A4_PCM/I2S_SDO
		pwmPin: -1,
	},
	"36": {
		pin:    223, // GPIO7_A7_UART3RX
		pwmPin: -1,
	},
	"37": {
		pin:    224, // GPIO7_B0_UART3TX
		pwmPin: -1,
	},
	"27": {
		pin:    233, // GPIO7_C1_I2C4_SDA
		pwmPin: -1,
	},
	"28": {
		pin:    234, // GPIO7_C2_I2C_SCL
		pwmPin: -1,
	},
	"33": {
		pin:    238, // GPIO7_C6_UART2RX_PWM2
		pwmPin: 0,
	},
	"32": {
		pin:    239, // GPIO7_C7_UART2TX_PWM3
		pwmPin: 1,
	},
	"26": {
		pin:    251, // GPIO8_A3_SPI2CSN1
		pwmPin: -1,
	},
	"3": {
		pin:    252, // GPIO8_A4_I2C1_SDA
		pwmPin: -1,
	},
	"5": {
		pin:    253, // GPIO8_A5_I2C1_SCL
		pwmPin: -1,
	},
	"23": {
		pin:    254, // GPIO8_A6_SPI2CLK
		pwmPin: -1,
	},
	"24": {
		pin:    255, // GPIO8_A7_SPI2CSN0
		pwmPin: -1,
	},
	"21": {
		pin:    256, // GPIO8_B0_SPI2RXD
		pwmPin: -1,
	},
	"19": {
		pin:    257, // GPIO8_B1_SPI2TXD
		pwmPin: -1,
	},
}
