package nanopi

import "gobot.io/x/gobot/v2/platforms/adaptors"

// pin definition for NanoPi NEO
// pins: A=0+Nr, C=64+Nr, G=192+Nr
var neoDigitalPinDefinitions = adaptors.DigitalPinDefinitions{
	"11": {Sysfs: 0, Cdev: adaptors.CdevPin{Chip: 0, Line: 0}},     // UART2_TX/GPIOA0
	"22": {Sysfs: 1, Cdev: adaptors.CdevPin{Chip: 0, Line: 1}},     // UART2_RX/GPIOA1
	"13": {Sysfs: 2, Cdev: adaptors.CdevPin{Chip: 0, Line: 2}},     // UART2_RTS/GPIOA2
	"15": {Sysfs: 3, Cdev: adaptors.CdevPin{Chip: 0, Line: 3}},     // UART2_CTS/GPIOA3
	"12": {Sysfs: 6, Cdev: adaptors.CdevPin{Chip: 0, Line: 6}},     // GPIOA6
	"19": {Sysfs: 64, Cdev: adaptors.CdevPin{Chip: 0, Line: 64}},   // SPI0_SDO/GPIOC0
	"21": {Sysfs: 65, Cdev: adaptors.CdevPin{Chip: 0, Line: 65}},   // SPI0_SDI/GPIOC1
	"23": {Sysfs: 66, Cdev: adaptors.CdevPin{Chip: 0, Line: 66}},   // SPI0_CLK/GPIOC2
	"24": {Sysfs: 67, Cdev: adaptors.CdevPin{Chip: 0, Line: 67}},   // SPI0_CS/GPIOC3
	"8":  {Sysfs: 198, Cdev: adaptors.CdevPin{Chip: 0, Line: 198}}, // UART1_TX/GPIOG6
	"10": {Sysfs: 199, Cdev: adaptors.CdevPin{Chip: 0, Line: 199}}, // UART1_RX/GPIOG7
	"16": {Sysfs: 200, Cdev: adaptors.CdevPin{Chip: 0, Line: 200}}, // UART1_RTS/GPIOG8
	"18": {Sysfs: 201, Cdev: adaptors.CdevPin{Chip: 0, Line: 201}}, // UART1_CTS/GPIOG9
	"7":  {Sysfs: 203, Cdev: adaptors.CdevPin{Chip: 0, Line: 203}}, // GPIOG11
}

var neoPWMPinDefinitions = adaptors.PWMPinDefinitions{
	// UART_RXD0, GPIOA5, PWM
	"PWM": {Dir: "/sys/devices/platform/soc/1c21400.pwm/pwm/", DirRegexp: "pwmchip[0]$", Channel: 0},
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	"thermal_zone0": {Path: "/sys/class/thermal/thermal_zone0/temp", R: true, W: false, BufLen: 7},
}
