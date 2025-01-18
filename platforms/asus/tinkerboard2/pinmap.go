package tinkerboard2

import "gobot.io/x/gobot/v2/platforms/adaptors"

// notes for character device
// pins: A=0+Nr, B=8+Nr, C=16+Nr, D=24+Nr
// tested: armbian Linux, OK: work as input and output, IN: work only as input
// MIPI ports itself not wired yet
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"26":    {Sysfs: 6, Cdev: adaptors.CdevPin{Chip: 0, Line: 6}},    // GPIO0_A6_PWM3A_IR - OK, PWM OK
	"7":     {Sysfs: 8, Cdev: adaptors.CdevPin{Chip: 0, Line: 8}},    // GPIO0_B0_TEST_CLKOUT2 - OK
	"21":    {Sysfs: 39, Cdev: adaptors.CdevPin{Chip: 1, Line: 7}},   // GPIO1_A7_SPI1_RXD_UART4_RXD - OK
	"19":    {Sysfs: 40, Cdev: adaptors.CdevPin{Chip: 1, Line: 8}},   // GPIO1_B0_SPI1_TXD_UART4_TXD - OK
	"23":    {Sysfs: 41, Cdev: adaptors.CdevPin{Chip: 1, Line: 9}},   // GPIO1_B1_SPI1_CLK - OK
	"24":    {Sysfs: 42, Cdev: adaptors.CdevPin{Chip: 1, Line: 10}},  // GPIO1_B2_SPI1_CSN - OK
	"DSI_1": {Sysfs: 48, Cdev: adaptors.CdevPin{Chip: 1, Line: 20}},  // GPIO1_C4_I2C8_SDA - ?
	"DSI_2": {Sysfs: 48, Cdev: adaptors.CdevPin{Chip: 1, Line: 21}},  // GPIO1_C5_I2C8_SCL - ?
	"27":    {Sysfs: 71, Cdev: adaptors.CdevPin{Chip: 2, Line: 7}},   // GPIO2_A7_I2C7_SDA - OK
	"28":    {Sysfs: 72, Cdev: adaptors.CdevPin{Chip: 2, Line: 8}},   // GPIO2_B0_I2C7_SCL - OK
	"3":     {Sysfs: 73, Cdev: adaptors.CdevPin{Chip: 2, Line: 9}},   // GPIO2_B1_I2C6_SDA - OK
	"5":     {Sysfs: 74, Cdev: adaptors.CdevPin{Chip: 2, Line: 10}},  // GPIO2_B2_I2C6_SCL - OK
	"CSI_3": {Sysfs: 75, Cdev: adaptors.CdevPin{Chip: 2, Line: 11}},  // GPIO2_B3_CSI_CLKOUT - ?
	"10":    {Sysfs: 80, Cdev: adaptors.CdevPin{Chip: 2, Line: 16}},  // GPIO2_C0_UART0_RXD - OK
	"8":     {Sysfs: 81, Cdev: adaptors.CdevPin{Chip: 2, Line: 17}},  // GPIO2_C1_UART0_TXD - OK
	"36":    {Sysfs: 82, Cdev: adaptors.CdevPin{Chip: 2, Line: 18}},  // GPIO2_C2_UART0_CTSN - OK
	"11":    {Sysfs: 83, Cdev: adaptors.CdevPin{Chip: 2, Line: 19}},  // GPIO2_C3_UART0_RTSN - OK
	"15":    {Sysfs: 84, Cdev: adaptors.CdevPin{Chip: 2, Line: 20}},  // GPIO2_C4_SPI5_RX - OK
	"13":    {Sysfs: 85, Cdev: adaptors.CdevPin{Chip: 2, Line: 21}},  // GPIO2_C5_SPI5_TX - OK
	"16":    {Sysfs: 86, Cdev: adaptors.CdevPin{Chip: 2, Line: 22}},  // GPIO2_C6_SPI5_CLK - OK
	"18":    {Sysfs: 87, Cdev: adaptors.CdevPin{Chip: 2, Line: 23}},  // GPIO2_C7_SPI5_CSN - OK
	"12":    {Sysfs: 120, Cdev: adaptors.CdevPin{Chip: 3, Line: 24}}, // GPIO3_D0_I2S0_SCLK - OK
	"35":    {Sysfs: 121, Cdev: adaptors.CdevPin{Chip: 3, Line: 25}}, // GPIO3_D1_I2S0_FS - OK (smooth digital behavior)
	"38":    {Sysfs: 123, Cdev: adaptors.CdevPin{Chip: 3, Line: 27}}, // GPIO3_D3_I2S0_SDI0 - OK
	"22":    {Sysfs: 124, Cdev: adaptors.CdevPin{Chip: 3, Line: 28}}, // GPIO3_D4_I2S0_SDO3 - OK
	"31":    {Sysfs: 125, Cdev: adaptors.CdevPin{Chip: 3, Line: 29}}, // GPIO3_D5_I2S0_SDO2 - OK
	"29":    {Sysfs: 126, Cdev: adaptors.CdevPin{Chip: 3, Line: 30}}, // GPIO3_D6_I2S0_SDO1 - OK
	"40":    {Sysfs: 127, Cdev: adaptors.CdevPin{Chip: 3, Line: 31}}, // GPIO3_D7_I2S0_SDO0 - OK
	"CSI_1": {Sysfs: 128, Cdev: adaptors.CdevPin{Chip: 4, Line: 1}},  // GPIO4_A1_I2C1_SDA -?
	"CSI_2": {Sysfs: 129, Cdev: adaptors.CdevPin{Chip: 4, Line: 2}},  // GPIO4_A2_I2C1_SCL -?
	"CSI_4": {Sysfs: 130, Cdev: adaptors.CdevPin{Chip: 4, Line: 3}},  // GPIO4_A3_CSI_GPIO -?
	"32":    {Sysfs: 146, Cdev: adaptors.CdevPin{Chip: 4, Line: 18}}, // GPIO4_C2_PWM0 - OK, PWM OK
	"J6_1":  {Sysfs: 147, Cdev: adaptors.CdevPin{Chip: 4, Line: 19}}, // GPIO4_C3_UART2_RX -?
	"J6_2":  {Sysfs: 148, Cdev: adaptors.CdevPin{Chip: 4, Line: 20}}, // GPIO4_C4_UART2_TX -?
	"37":    {Sysfs: 149, Cdev: adaptors.CdevPin{Chip: 4, Line: 21}}, // GPIO4_C5_SPDIF_TX - OK
	"33":    {Sysfs: 150, Cdev: adaptors.CdevPin{Chip: 4, Line: 22}}, // GPIO4_C6_PWM1 - OK, PWM OK
}

var pwmPinDefinitions = adaptors.PWMPinDefinitions{
	// needs to be enabled by device tree (pwm0)
	"32": {Dir: "/sys/devices/platform/ff420000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3]$", Channel: 0},
	// needs to be enabled by device tree (pwm1)
	"33": {Dir: "/sys/devices/platform/ff420010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3]$", Channel: 0},
	// needs to be enabled by device tree (pwm3)
	"26": {Dir: "/sys/devices/platform/ff420030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3]$", Channel: 0},
}

// analog pins are the same as for tinkerboard
