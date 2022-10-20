package tinkerboard

var gpioPinDefinitions = map[string]int{
	"7":  17,  // GPIO0_C1_CLKOUT
	"10": 160, // GPIO5_B0_UART1RX
	"8":  161, // GPIO5_B1_UART1TX
	"16": 162, // GPIO5_B2_UART1CTSN
	"18": 163, // GPIO5_B3_UART1RTSN
	"11": 164, // GPIO5_B4_SPI0CLK_UART4CTSN
	"29": 165, // GPIO5_B5_SPI0CSN_UART4RTSN
	"13": 166, // GPIO5_B6_SPI0_TXD_UART4TX
	"15": 167, // GPIO5_B7_SPI0_RXD_UART4RX
	"31": 168, // GPIO5_C0_SPI0CSN1
	"22": 171, // GPIO5_C3
	"12": 184, // GPIO6_A0_PCM/I2S_CLK
	"35": 185, // GPIO6_A1_PCM/I2S_FS
	"38": 187, // GPIO6_A3_PCM/I2S_SDI
	"40": 188, // GPIO6_A4_PCM/I2S_SDO
	"36": 223, // GPIO7_A7_UART3RX
	"37": 224, // GPIO7_B0_UART3TX
	"27": 233, // GPIO7_C1_I2C4_SDA
	"28": 234, // GPIO7_C2_I2C_SCL
	"33": 238, // GPIO7_C6_UART2RX_PWM2
	"32": 239, // GPIO7_C7_UART2TX_PWM3
	"26": 251, // GPIO8_A3_SPI2CSN1
	"3":  252, // GPIO8_A4_I2C1_SDA
	"5":  253, // GPIO8_A5_I2C1_SCL
	"23": 254, // GPIO8_A6_SPI2CLK
	"24": 255, // GPIO8_A7_SPI2CSN0
	"21": 256, // GPIO8_B0_SPI2RXD
	"19": 257, // GPIO8_B1_SPI2TXD
}

var pwmPinDefinitions = map[string]pwmPinDefinition{
	// GPIO7_C6_UART2RX_PWM2
	"33": pwmPinDefinition{dir: "/sys/devices/platform/ff680020.pwm/pwm/", dirRegexp: "pwmchip2$", channel: 0},
	// GPIO7_C7_UART2TX_PWM3
	"32": pwmPinDefinition{dir: "/sys/devices/platform/ff680030.pwm/pwm/", dirRegexp: "pwmchip[2|3]$", channel: 0},
}
