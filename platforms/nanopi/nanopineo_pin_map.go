package nanopi

// pin definition for NanoPi NEO
// pins: A=0+Nr, C=64+Nr, G=192+Nr
var neoGpioPins = map[string]gpioPinDefinition{
	"11": gpioPinDefinition{sysfs: 0, cdev: cdevPin{chip: 0, line: 0}},     // UART2_TX/GPIOA0
	"22": gpioPinDefinition{sysfs: 1, cdev: cdevPin{chip: 0, line: 1}},     // UART2_RX/GPIOA1
	"13": gpioPinDefinition{sysfs: 2, cdev: cdevPin{chip: 0, line: 2}},     // UART2_RTS/GPIOA2
	"15": gpioPinDefinition{sysfs: 3, cdev: cdevPin{chip: 0, line: 3}},     // UART2_CTS/GPIOA3
	"12": gpioPinDefinition{sysfs: 6, cdev: cdevPin{chip: 0, line: 6}},     // GPIOA6
	"19": gpioPinDefinition{sysfs: 64, cdev: cdevPin{chip: 0, line: 64}},   // SPI0_MOSI/GPIOC0
	"21": gpioPinDefinition{sysfs: 65, cdev: cdevPin{chip: 0, line: 65}},   // SPI0_MISO/GPIOC1
	"23": gpioPinDefinition{sysfs: 66, cdev: cdevPin{chip: 0, line: 66}},   // SPI0_CLK/GPIOC2
	"24": gpioPinDefinition{sysfs: 67, cdev: cdevPin{chip: 0, line: 67}},   // SPI0_CS/GPIOC3
	"8":  gpioPinDefinition{sysfs: 198, cdev: cdevPin{chip: 0, line: 198}}, // UART1_TX/GPIOG6
	"10": gpioPinDefinition{sysfs: 199, cdev: cdevPin{chip: 0, line: 199}}, // UART1_RX/GPIOG7
	"16": gpioPinDefinition{sysfs: 200, cdev: cdevPin{chip: 0, line: 200}}, // UART1_RTS/GPIOG8
	"18": gpioPinDefinition{sysfs: 201, cdev: cdevPin{chip: 0, line: 201}}, // UART1_CTS/GPIOG9
	"7":  gpioPinDefinition{sysfs: 203, cdev: cdevPin{chip: 0, line: 203}}, // GPIOG11
}

var neoPwmPins = map[string]pwmPinDefinition{
	// UART_RXD0, GPIOA5, PWM
	"PWM": pwmPinDefinition{dir: "/sys/devices/platform/soc/1c21400.pwm/pwm/", dirRegexp: "pwmchip[0]$", channel: 0},
}
