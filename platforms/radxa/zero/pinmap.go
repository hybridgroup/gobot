package zero

import "gobot.io/x/gobot/v2/platforms/adaptors"

// tested with cdev on a Zero V1.51 board: dietPi Linux, OK: works, ?: unknown, NOK: not working
// IN: works only as input, PU: if used as input, external pullup resistor needed
//
//nolint:lll // ok here
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"8":   {Sysfs: 412, Cdev: adaptors.CdevPin{Chip: 1, Line: 0}},  // GPIOAO_0_UART_AO_A_TXD - ?
	"10":  {Sysfs: 413, Cdev: adaptors.CdevPin{Chip: 1, Line: 1}},  // GPIOAO_1_UART_AO_A_RXD - OK
	"11":  {Sysfs: 414, Cdev: adaptors.CdevPin{Chip: 1, Line: 2}},  // GPIOAO_2_I2C_AO_M0_SCL_UART_AO_B_TX_I2C_AO_S0_SCL - OK
	"28":  {Sysfs: 414, Cdev: adaptors.CdevPin{Chip: 1, Line: 2}},  // GPIOAO_2_I2C_AO_M0_SCL_UART_AO_B_TX_I2C_AO_S0_SCL - OK
	"7":   {Sysfs: 415, Cdev: adaptors.CdevPin{Chip: 1, Line: 3}},  // GPIOAO_3_I2C_AO_M0_SDA_UART_AO_B_RX_I2C_AO_S0_SDA - OK
	"27":  {Sysfs: 415, Cdev: adaptors.CdevPin{Chip: 1, Line: 3}},  // GPIOAO_3_I2C_AO_M0_SDA_UART_AO_B_RX_I2C_AO_S0_SDA - OK
	"32":  {Sysfs: 416, Cdev: adaptors.CdevPin{Chip: 1, Line: 4}},  // GPIOAO_4_PWMAO_C - OK
	"35":  {Sysfs: 420, Cdev: adaptors.CdevPin{Chip: 1, Line: 8}},  // GPIOAO_8_UART_AO_B_TX - OK
	"37":  {Sysfs: 421, Cdev: adaptors.CdevPin{Chip: 1, Line: 9}},  // GPIOAO_9_UART_AO_B_RX - OK
	"LED": {Sysfs: 422, Cdev: adaptors.CdevPin{Chip: 1, Line: 10}}, // GPIOAO_10_PWMAO_D - Wired to LED besides USB-C
	"40":  {Sysfs: 423, Cdev: adaptors.CdevPin{Chip: 1, Line: 11}}, // GPIOAO_11_PWMAO_A - OK
	"19":  {Sysfs: 447, Cdev: adaptors.CdevPin{Chip: 0, Line: 20}}, // GPIOH_4_UART_EE_C_RTS_SPI_B_MOSI - OK
	"21":  {Sysfs: 448, Cdev: adaptors.CdevPin{Chip: 0, Line: 21}}, // GPIOH_5_UART_EE_C_CTS_SPI_B_MISO_PWM_F - OK
	"24":  {Sysfs: 449, Cdev: adaptors.CdevPin{Chip: 0, Line: 22}}, // GPIOH_6_UART_EE_C_RX_SPI_B_SS0_I2C_EE_M1_SDA - OK
	"23":  {Sysfs: 450, Cdev: adaptors.CdevPin{Chip: 0, Line: 23}}, // GPIOH_7_UART_EE_C_TX_SPI_B_SCLK_I2C_EE_M1_SCL - OK
	"36":  {Sysfs: 451, Cdev: adaptors.CdevPin{Chip: 0, Line: 24}}, // GPIOH_8 - OK
	"22":  {Sysfs: 475, Cdev: adaptors.CdevPin{Chip: 0, Line: 48}}, // GPIOC_7 - OK
	"3":   {Sysfs: 490, Cdev: adaptors.CdevPin{Chip: 0, Line: 63}}, // GPIOA_14_I2C_EE_M3_SDA - OK (i2c-3)
	"5":   {Sysfs: 491, Cdev: adaptors.CdevPin{Chip: 0, Line: 64}}, // GPIOA_15_I2C_EE_M3_SCL - OK (i2-c3)
	"18":  {Sysfs: 500, Cdev: adaptors.CdevPin{Chip: 0, Line: 73}}, // GPIOX_8_SPI_A_MOSI_PWM_C_TDMA_D1 - OK
	"12":  {Sysfs: 501, Cdev: adaptors.CdevPin{Chip: 0, Line: 74}}, // GPIOX_9_SPI_A_MISO_TDMA_D0 - OK
	"16":  {Sysfs: 502, Cdev: adaptors.CdevPin{Chip: 0, Line: 75}}, // GPIOX_10_SPI_A_SS0_I2C_EE_M1_SDA_TDMA_FS - OK
	"13":  {Sysfs: 503, Cdev: adaptors.CdevPin{Chip: 0, Line: 76}}, // GPIOX_11_SPI_A_SCLK_I2C_EE_M1_SCL_TDMA_SCLK - OK
}

var pwmPinDefinitions = adaptors.PWMPinDefinitions{
	// enabled by default, but pin seems to be not really enabled on pwm usage, channel 1 used by "regulator-vddcpu"
	"32": {Dir: "/sys/devices/platform/soc/ff800000.bus/ff802000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4]$", Channel: 0},
	"40": {Dir: "/sys/devices/platform/soc/ff800000.bus/ff807000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4]$", Channel: 0},
	"18": {Dir: "/sys/devices/platform/soc/ffd00000.bus/ffd1a000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4]$", Channel: 0},
	// enabled by default, but pin seems to be not really enabled on pwm usage, channel 0 used by "wifi32k"
	"21": {Dir: "/sys/devices/platform/soc/ffd00000.bus/ffd19000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4]$", Channel: 1},
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	// names equals /sys/class/thermal/thermal_zone*/hwmon*/name
	"cpu_thermal": {Path: "/sys/class/thermal/thermal_zone0/temp", W: false, ReadBufLen: 7},
	"ddr_thermal": {Path: "/sys/class/thermal/thermal_zone1/temp", W: false, ReadBufLen: 7},
	// 2 channel 12-bit SAR ADC, 0..4095, so 4 characters to read
	"15": {
		Path: "/sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage1_raw", W: false, ReadBufLen: 4,
	},
	"15_mean": {
		Path: "/sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage1_mean_raw", W: false, ReadBufLen: 4,
	},
	"26": {
		Path: "/sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage2_raw", W: false, ReadBufLen: 4,
	},
	"26_mean": {
		Path: "/sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage2_mean_raw", W: false, ReadBufLen: 4,
	},
}
