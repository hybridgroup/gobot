package tinkerboard

import "gobot.io/x/gobot/v2/platforms/adaptors"

// notes for character device
// pins: A=0+Nr, B=8+Nr, C=16+Nr
// tested: armbian Linux, OK: work as input and output, IN: work only as input
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"7":  {Sysfs: 17, Cdev: adaptors.CdevPin{Chip: 0, Line: 17}},  // GPIO0_C1_CLKOUT - OK
	"10": {Sysfs: 160, Cdev: adaptors.CdevPin{Chip: 5, Line: 8}},  // GPIO5_B0_UART1RX - IN, initial 1
	"8":  {Sysfs: 161, Cdev: adaptors.CdevPin{Chip: 5, Line: 9}},  // GPIO5_B1_UART1TX - NO, initial 1
	"16": {Sysfs: 162, Cdev: adaptors.CdevPin{Chip: 5, Line: 10}}, // GPIO5_B2_UART1CTSN - NO, initial 0
	"18": {Sysfs: 163, Cdev: adaptors.CdevPin{Chip: 5, Line: 11}}, // GPIO5_B3_UART1RTSN - NO, initial 0
	"11": {Sysfs: 164, Cdev: adaptors.CdevPin{Chip: 5, Line: 12}}, // GPIO5_B4_SPI0CLK_UART4CTSN - NO, initial 0
	"29": {Sysfs: 165, Cdev: adaptors.CdevPin{Chip: 5, Line: 13}}, // GPIO5_B5_SPI0CSN_UART4RTSN - NO, initial 0
	"13": {Sysfs: 166, Cdev: adaptors.CdevPin{Chip: 5, Line: 14}}, // GPIO5_B6_SPI0_TXD_UART4TX - NO, initial 1
	"15": {Sysfs: 167, Cdev: adaptors.CdevPin{Chip: 5, Line: 15}}, // GPIO5_B7_SPI0_RXD_UART4RX - IN, initial 1
	"31": {Sysfs: 168, Cdev: adaptors.CdevPin{Chip: 5, Line: 16}}, // GPIO5_C0_SPI0CSN1 - OK if SPI0 off
	"22": {Sysfs: 171, Cdev: adaptors.CdevPin{Chip: 5, Line: 19}}, // GPIO5_C3 - OK
	"12": {Sysfs: 184, Cdev: adaptors.CdevPin{Chip: 6, Line: 0}},  // GPIO6_A0_PCM/I2S_CLK - NO, initial 1
	"35": {Sysfs: 185, Cdev: adaptors.CdevPin{Chip: 6, Line: 1}},  // GPIO6_A1_PCM/I2S_FS - NO, initial 0
	"38": {Sysfs: 187, Cdev: adaptors.CdevPin{Chip: 6, Line: 3}},  // GPIO6_A3_PCM/I2S_SDI - IN, initial 1
	"40": {Sysfs: 188, Cdev: adaptors.CdevPin{Chip: 6, Line: 4}},  // GPIO6_A4_PCM/I2S_SDO - NO, initial 0
	"36": {Sysfs: 223, Cdev: adaptors.CdevPin{Chip: 7, Line: 7}},  // GPIO7_A7_UART3RX - IN, initial 1
	"37": {Sysfs: 224, Cdev: adaptors.CdevPin{Chip: 7, Line: 8}},  // GPIO7_B0_UART3TX - NO, initial 1
	"27": {Sysfs: 233, Cdev: adaptors.CdevPin{Chip: 7, Line: 17}}, // GPIO7_C1_I2C4_SDA - OK if I2C4 off
	"28": {Sysfs: 234, Cdev: adaptors.CdevPin{Chip: 7, Line: 18}}, // GPIO7_C2_I2C_SCL - OK if I2C4 off
	"33": {Sysfs: 238, Cdev: adaptors.CdevPin{Chip: 7, Line: 22}}, // GPIO7_C6_UART2RX_PWM2 - IN, initial 1
	"32": {Sysfs: 239, Cdev: adaptors.CdevPin{Chip: 7, Line: 23}}, // GPIO7_C7_UART2TX_PWM3 - NO, initial 1
	"26": {Sysfs: 251, Cdev: adaptors.CdevPin{Chip: 8, Line: 3}},  // GPIO8_A3_SPI2CSN1 - OK if SPI2 off
	"3":  {Sysfs: 252, Cdev: adaptors.CdevPin{Chip: 8, Line: 4}},  // GPIO8_A4_I2C1_SDA - OK if I2C1 off
	"5":  {Sysfs: 253, Cdev: adaptors.CdevPin{Chip: 8, Line: 5}},  // GPIO8_A5_I2C1_SCL - OK if I2C1 off
	"23": {Sysfs: 254, Cdev: adaptors.CdevPin{Chip: 8, Line: 6}},  // GPIO8_A6_SPI2CLK - OK if SPI2 off
	"24": {Sysfs: 255, Cdev: adaptors.CdevPin{Chip: 8, Line: 7}},  // GPIO8_A7_SPI2CSN0 - OK if SPI2 off
	"21": {Sysfs: 256, Cdev: adaptors.CdevPin{Chip: 8, Line: 8}},  // GPIO8_B0_SPI2RXD - OK if SPI2 off
	"19": {Sysfs: 257, Cdev: adaptors.CdevPin{Chip: 8, Line: 9}},  // GPIO8_B1_SPI2TXD - OK if SPI2 off
}

var pwmPinDefinitions = adaptors.PWMPinDefinitions{
	// GPIO7_C6_UART2RX_PWM2
	"33": {Dir: "/sys/devices/platform/ff680020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2]$", Channel: 0},
	// GPIO7_C7_UART2TX_PWM3
	"32": {Dir: "/sys/devices/platform/ff680030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3]$", Channel: 0},
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	"thermal_zone0": {Path: "/sys/class/thermal/thermal_zone0/temp", R: true, W: false, BufLen: 7},
	"thermal_zone1": {Path: "/sys/class/thermal/thermal_zone1/temp", R: true, W: false, BufLen: 7},
}
