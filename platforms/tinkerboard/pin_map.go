package tinkerboard

type cdevPin struct {
	chip uint8
	line uint8
}

type gpioPinDefinition struct {
	sysfs int
	cdev  cdevPin
}

// notes for character device
// pins: A=0+Nr, B=8+Nr, C=16+Nr
// tested: armbian Linux, OK: work as input and output, IN: work only as input
var gpioPinDefinitions = map[string]gpioPinDefinition{
	"7":  {sysfs: 17, cdev: cdevPin{chip: 0, line: 17}},  // GPIO0_C1_CLKOUT - OK
	"10": {sysfs: 160, cdev: cdevPin{chip: 5, line: 8}},  // GPIO5_B0_UART1RX - IN, initial 1
	"8":  {sysfs: 161, cdev: cdevPin{chip: 5, line: 9}},  // GPIO5_B1_UART1TX - NO, initial 1
	"16": {sysfs: 162, cdev: cdevPin{chip: 5, line: 10}}, // GPIO5_B2_UART1CTSN - NO, initial 0
	"18": {sysfs: 163, cdev: cdevPin{chip: 5, line: 11}}, // GPIO5_B3_UART1RTSN - NO, initial 0
	"11": {sysfs: 164, cdev: cdevPin{chip: 5, line: 12}}, // GPIO5_B4_SPI0CLK_UART4CTSN - NO, initial 0
	"29": {sysfs: 165, cdev: cdevPin{chip: 5, line: 13}}, // GPIO5_B5_SPI0CSN_UART4RTSN - NO, initial 0
	"13": {sysfs: 166, cdev: cdevPin{chip: 5, line: 14}}, // GPIO5_B6_SPI0_TXD_UART4TX - NO, initial 1
	"15": {sysfs: 167, cdev: cdevPin{chip: 5, line: 15}}, // GPIO5_B7_SPI0_RXD_UART4RX - IN, initial 1
	"31": {sysfs: 168, cdev: cdevPin{chip: 5, line: 16}}, // GPIO5_C0_SPI0CSN1 - OK if SPI0 off
	"22": {sysfs: 171, cdev: cdevPin{chip: 5, line: 19}}, // GPIO5_C3 - OK
	"12": {sysfs: 184, cdev: cdevPin{chip: 6, line: 0}},  // GPIO6_A0_PCM/I2S_CLK - NO, initial 1
	"35": {sysfs: 185, cdev: cdevPin{chip: 6, line: 1}},  // GPIO6_A1_PCM/I2S_FS - NO, initial 0
	"38": {sysfs: 187, cdev: cdevPin{chip: 6, line: 3}},  // GPIO6_A3_PCM/I2S_SDI - IN, initial 1
	"40": {sysfs: 188, cdev: cdevPin{chip: 6, line: 4}},  // GPIO6_A4_PCM/I2S_SDO - NO, initial 0
	"36": {sysfs: 223, cdev: cdevPin{chip: 7, line: 7}},  // GPIO7_A7_UART3RX - IN, initial 1
	"37": {sysfs: 224, cdev: cdevPin{chip: 7, line: 8}},  // GPIO7_B0_UART3TX - NO, initial 1
	"27": {sysfs: 233, cdev: cdevPin{chip: 7, line: 17}}, // GPIO7_C1_I2C4_SDA - OK if I2C4 off
	"28": {sysfs: 234, cdev: cdevPin{chip: 7, line: 18}}, // GPIO7_C2_I2C_SCL - OK if I2C4 off
	"33": {sysfs: 238, cdev: cdevPin{chip: 7, line: 22}}, // GPIO7_C6_UART2RX_PWM2 - IN, initial 1
	"32": {sysfs: 239, cdev: cdevPin{chip: 7, line: 23}}, // GPIO7_C7_UART2TX_PWM3 - NO, initial 1
	"26": {sysfs: 251, cdev: cdevPin{chip: 8, line: 3}},  // GPIO8_A3_SPI2CSN1 - OK if SPI2 off
	"3":  {sysfs: 252, cdev: cdevPin{chip: 8, line: 4}},  // GPIO8_A4_I2C1_SDA - OK if I2C1 off
	"5":  {sysfs: 253, cdev: cdevPin{chip: 8, line: 5}},  // GPIO8_A5_I2C1_SCL - OK if I2C1 off
	"23": {sysfs: 254, cdev: cdevPin{chip: 8, line: 6}},  // GPIO8_A6_SPI2CLK - OK if SPI2 off
	"24": {sysfs: 255, cdev: cdevPin{chip: 8, line: 7}},  // GPIO8_A7_SPI2CSN0 - OK if SPI2 off
	"21": {sysfs: 256, cdev: cdevPin{chip: 8, line: 8}},  // GPIO8_B0_SPI2RXD - OK if SPI2 off
	"19": {sysfs: 257, cdev: cdevPin{chip: 8, line: 9}},  // GPIO8_B1_SPI2TXD - OK if SPI2 off
}

var pwmPinDefinitions = map[string]pwmPinDefinition{
	// GPIO7_C6_UART2RX_PWM2
	"33": {dir: "/sys/devices/platform/ff680020.pwm/pwm/", dirRegexp: "pwmchip[0|1|2]$", channel: 0},
	// GPIO7_C7_UART2TX_PWM3
	"32": {dir: "/sys/devices/platform/ff680030.pwm/pwm/", dirRegexp: "pwmchip[0|1|2|3]$", channel: 0},
}
