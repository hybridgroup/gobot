package nanopct6

import "gobot.io/x/gobot/v2/platforms/adaptors"

// notes for character device
// sysfs: Chip*32 + (A=0, B=8, C=16) + Nr
// tested with cdev on a NanoPC-T6 2301 board: armbian Linux, OK: works, ?: unknown, NOK: not working
// IN: works only as input, PU: if used as input, external pullup resistor needed
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"10":      {Sysfs: 20, Cdev: adaptors.CdevPin{Chip: 0, Line: 20}},  // GPIO0_C4_UART0_RX_M0 - OK
	"8":       {Sysfs: 21, Cdev: adaptors.CdevPin{Chip: 0, Line: 21}},  // GPIO0_C5_UART0_TX_M0_PWM4_M0 - OK
	"32":      {Sysfs: 22, Cdev: adaptors.CdevPin{Chip: 0, Line: 22}},  // GPIO0_C6_PWM5_M1 - OK
	"27":      {Sysfs: 32, Cdev: adaptors.CdevPin{Chip: 1, Line: 0}},   // GPIO1_A0_UART6_RX_M1 - OK
	"28":      {Sysfs: 33, Cdev: adaptors.CdevPin{Chip: 1, Line: 1}},   // GPIO1_A1_UART6_TX_M1 - OK (UP)
	"15":      {Sysfs: 39, Cdev: adaptors.CdevPin{Chip: 1, Line: 7}},   // GPIO1_A7 - OK
	"26":      {Sysfs: 40, Cdev: adaptors.CdevPin{Chip: 1, Line: 8}},   // GPIO1_B0 - OK
	"21":      {Sysfs: 41, Cdev: adaptors.CdevPin{Chip: 1, Line: 9}},   // GPIO1_B1_SPI0_MISO_M2 - OK (UP)
	"19":      {Sysfs: 42, Cdev: adaptors.CdevPin{Chip: 1, Line: 10}},  // GPIO1_B2_SPI0_MOSI_M2_UART4_RX_M2 - OK
	"23":      {Sysfs: 43, Cdev: adaptors.CdevPin{Chip: 1, Line: 11}},  // GPIO1_B3_SPI0_CLK_M2_UART4_TX_M2 - OK
	"24":      {Sysfs: 44, Cdev: adaptors.CdevPin{Chip: 1, Line: 12}},  // GPIO1_B4_SPI0_CS0_M2_UART7_RX_M2 - OK
	"22":      {Sysfs: 45, Cdev: adaptors.CdevPin{Chip: 1, Line: 13}},  // GPIO1_B5_SPI0_CS1_M0_UART7_TX_M2 - OK
	"5":       {Sysfs: 62, Cdev: adaptors.CdevPin{Chip: 1, Line: 30}},  // GPIO1_D6_I2C8_SCL_M2 - OK
	"3":       {Sysfs: 63, Cdev: adaptors.CdevPin{Chip: 1, Line: 31}},  // GPIO1_D7_I2C8_SDA_M2 - OK
	"CSI1_11": {Sysfs: 81, Cdev: adaptors.CdevPin{Chip: 2, Line: 17}},  // GPIO2_C1 - ?
	"CSI1_12": {Sysfs: 82, Cdev: adaptors.CdevPin{Chip: 2, Line: 18}},  // GPIO2_C2 - ?
	"35":      {Sysfs: 96, Cdev: adaptors.CdevPin{Chip: 3, Line: 0}},   // GPIO3_A0_SPI4_MISO_M1_I2S3_MCLK_PWM10_M0 - OK
	"38":      {Sysfs: 97, Cdev: adaptors.CdevPin{Chip: 3, Line: 1}},   // GPIO3_A1_SPI4_MOSI_M1_I2S3_SCLK - OK
	"40":      {Sysfs: 98, Cdev: adaptors.CdevPin{Chip: 3, Line: 2}},   // GPIO3_A2_SPI4_CLK_M1_UART8_TX_M1_I2S3_LRCK - OK
	"36":      {Sysfs: 99, Cdev: adaptors.CdevPin{Chip: 3, Line: 3}},   // GPIO3_A3_SPI4_CS0_M1_UART8_RX_M1_I2S3_SDO - OK
	"37":      {Sysfs: 100, Cdev: adaptors.CdevPin{Chip: 3, Line: 4}},  // GPIO3_A4_SPI4_CS1_M1_I2S3_SDI - OK
	"DSI0_12": {Sysfs: 102, Cdev: adaptors.CdevPin{Chip: 3, Line: 6}},  // GPIO3_A6 - ?
	"33":      {Sysfs: 104, Cdev: adaptors.CdevPin{Chip: 3, Line: 8}},  // GPIO3_B0_PWM9_M0 - OK
	"DSI0_10": {Sysfs: 105, Cdev: adaptors.CdevPin{Chip: 3, Line: 9}},  // GPIO3_B1_PWM2_M1 - ?
	"7":       {Sysfs: 106, Cdev: adaptors.CdevPin{Chip: 3, Line: 10}}, // GPIO3_B2_I2S2_SDI_M1 - OK
	"16":      {Sysfs: 107, Cdev: adaptors.CdevPin{Chip: 3, Line: 11}}, // GPIO3_B3_I2S2_SDO_M1 - OK
	"18":      {Sysfs: 108, Cdev: adaptors.CdevPin{Chip: 3, Line: 12}}, // GPIO3_B4_I2S2_MCLK_M1 - OK
	"29":      {Sysfs: 109, Cdev: adaptors.CdevPin{Chip: 3, Line: 13}}, // GPIO3_B5_UART3_TX_M1_I2S2_SCLK_M1_PWM12_M0 - OK
	"31":      {Sysfs: 110, Cdev: adaptors.CdevPin{Chip: 3, Line: 14}}, // GPIO3_B6_UART3_RX_M1_I2S2_LRCK_M1_PWM13_M0 - OK
	"12":      {Sysfs: 111, Cdev: adaptors.CdevPin{Chip: 3, Line: 15}}, // GPIO3_B7 - OK (UP)
	"DSI0_8":  {Sysfs: 112, Cdev: adaptors.CdevPin{Chip: 3, Line: 16}}, // GPIO3_C0 - ?
	"DSI0_14": {Sysfs: 113, Cdev: adaptors.CdevPin{Chip: 3, Line: 17}}, // GPIO3_C1 - ?
	"11":      {Sysfs: 114, Cdev: adaptors.CdevPin{Chip: 3, Line: 18}}, // GPIO3_C2_PWM14_M0 - OK
	"13":      {Sysfs: 115, Cdev: adaptors.CdevPin{Chip: 3, Line: 19}}, // GPIO3_C3_PWM15_IR_M0 - OK
	"DSI1_10": {Sysfs: 125, Cdev: adaptors.CdevPin{Chip: 3, Line: 29}}, // GPIO3_D5_PWM11_M3 - ?
	"DSI1_8":  {Sysfs: 128, Cdev: adaptors.CdevPin{Chip: 4, Line: 0}},  // GPIO4_A0 - ?
	"DSI1_14": {Sysfs: 129, Cdev: adaptors.CdevPin{Chip: 4, Line: 1}},  // GPIO4_A1 - ?
	"DSI1_12": {Sysfs: 131, Cdev: adaptors.CdevPin{Chip: 4, Line: 3}},  // GPIO4_A3 - ?
	"CSI0_11": {Sysfs: 148, Cdev: adaptors.CdevPin{Chip: 4, Line: 20}}, // GPIO4_C4 - ?
	"CSI0_12": {Sysfs: 149, Cdev: adaptors.CdevPin{Chip: 4, Line: 21}}, // GPIO4_C5 - ?
}

var pwmPinDefinitions = adaptors.PWMPinDefinitions{
	// needs to be enabled by DT-overlay pwm2-m1 (pwm2 = "/pwm@fd8b0020";)
	"DSI0_10": {Dir: "/sys/devices/platform/fd8b0020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm4-m0 (pwm4 = "/pwm@febd0000";)
	"8": {Dir: "/sys/devices/platform/febd0000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm5-m1 (pwm5 = "/pwm@febd0010";)
	"32": {Dir: "/sys/devices/platform/febd0010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm9-m0 (pwm9 = "/pwm@febe0010";)
	"33": {Dir: "/sys/devices/platform/febe0010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm10-m0 (pwm10 = "/pwm@febe0020";)
	"35": {Dir: "/sys/devices/platform/febe0020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm11-m3 (pwm11 = "/pwm@febe0030";)
	"DSI1_10": {Dir: "/sys/devices/platform/febe0030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm12-m0 (pwm12 = "/pwm@febf0000";)
	"29": {Dir: "/sys/devices/platform/febf0000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm13-m0  (pwm13 = "/pwm@febf0010";)
	"31": {Dir: "/sys/devices/platform/febf0010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm14-m0 (pwm14 = "/pwm@febf0020";)
	"11": {Dir: "/sys/devices/platform/febf0020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm15-m0 (pwm15 = "/pwm@febf0030";)
	"13": {Dir: "/sys/devices/platform/febf0030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	// names equals /sys/class/thermal/thermal_zone*/hwmon*/name
	"soc_thermal":        {Path: "/sys/class/thermal/thermal_zone0/temp", W: false, ReadBufLen: 7},
	"bigcore0_thermal":   {Path: "/sys/class/thermal/thermal_zone1/temp", W: false, ReadBufLen: 7},
	"bigcore1_thermal":   {Path: "/sys/class/thermal/thermal_zone2/temp", W: false, ReadBufLen: 7},
	"littlecore_thermal": {Path: "/sys/class/thermal/thermal_zone3/temp", W: false, ReadBufLen: 7},
	"center_thermal":     {Path: "/sys/class/thermal/thermal_zone4/temp", W: false, ReadBufLen: 7},
	"gpu_thermal":        {Path: "/sys/class/thermal/thermal_zone5/temp", W: false, ReadBufLen: 7},
	"npu_thermal":        {Path: "/sys/class/thermal/thermal_zone6/temp", W: false, ReadBufLen: 7},
}
