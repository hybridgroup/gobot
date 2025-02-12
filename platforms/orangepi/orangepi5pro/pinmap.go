package orangepi5pro

import "gobot.io/x/gobot/v2/platforms/adaptors"

// notes for character device
// sysfs: Chip*32 + (A=0, B=8, C=16, D=24) + Nr
// tested with cdev on a OrangePi 5 Pro V1.2 board: armbian Linux, OK: works, ?: unknown, NOK: not working
// IN: works only as input, PU: if used as input, external pullup resistor needed
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"8":  {Sysfs: 13, Cdev: adaptors.CdevPin{Chip: 0, Line: 13}},  // GPIOO_B5_UART2_TX_MO - OK
	"10": {Sysfs: 14, Cdev: adaptors.CdevPin{Chip: 0, Line: 14}},  // GPIOO_B6_UART2_RX_MO - OK
	"18": {Sysfs: 32, Cdev: adaptors.CdevPin{Chip: 1, Line: 0}},   // GPIO1_A0_UART6_RX_M1 - OK
	"16": {Sysfs: 33, Cdev: adaptors.CdevPin{Chip: 1, Line: 1}},   // GPIO1_A1_UART6_TX_M1 - PU
	"27": {Sysfs: 34, Cdev: adaptors.CdevPin{Chip: 1, Line: 2}},   // GPIO1_A2_I2C4_SDA_M3_PWM0_M2_SPI4_CLK_M2 - PU
	"28": {Sysfs: 35, Cdev: adaptors.CdevPin{Chip: 1, Line: 3}},   // GPIO1_A3_I2C4_SCL_M3_PWM1_M2_SPI4_CSO_M2 - OK
	"29": {Sysfs: 36, Cdev: adaptors.CdevPin{Chip: 1, Line: 4}},   // GPIO1_A4 - OK
	"31": {Sysfs: 38, Cdev: adaptors.CdevPin{Chip: 1, Line: 6}},   // GPIO1_A6 - OK
	"12": {Sysfs: 39, Cdev: adaptors.CdevPin{Chip: 1, Line: 7}},   // GPIO1_A7_PWM3_IR_M3 - OK
	"22": {Sysfs: 40, Cdev: adaptors.CdevPin{Chip: 1, Line: 8}},   // GPIO1_B0 - OK
	"21": {Sysfs: 41, Cdev: adaptors.CdevPin{Chip: 1, Line: 9}},   // GPIO1_B1_SPI0_MISO_M2 - OK
	"19": {Sysfs: 42, Cdev: adaptors.CdevPin{Chip: 1, Line: 10}},  // GPIO1_B2_SPI0_MOSI_M2_UART4_RX_M2 - PU
	"23": {Sysfs: 43, Cdev: adaptors.CdevPin{Chip: 1, Line: 11}},  // GPIO1_B3_SPI0_CLK_M2_UART4_TX_M2 - OK
	"24": {Sysfs: 44, Cdev: adaptors.CdevPin{Chip: 1, Line: 12}},  // GPIO1_B4_SPI0_CSO_M2_UART7_RX_M2 - OK
	"26": {Sysfs: 45, Cdev: adaptors.CdevPin{Chip: 1, Line: 13}},  // GPIO1_B5_SPI0_CS1_M2_UART7_TX_M2 - OK
	"15": {Sysfs: 46, Cdev: adaptors.CdevPin{Chip: 1, Line: 14}},  // GPIO1_B6_I2C5_SCL_M3_UART1_TX_M1 - OK
	"7":  {Sysfs: 47, Cdev: adaptors.CdevPin{Chip: 1, Line: 15}},  // GPIO1_B7_PWM13_M2_I2C5_SDA_M3_UART1_RX_M1 - OK
	"5":  {Sysfs: 58, Cdev: adaptors.CdevPin{Chip: 1, Line: 26}},  // GPIO1_D2_I2C1_SCL_M4_UART4_TX_M0 - OK
	"3":  {Sysfs: 59, Cdev: adaptors.CdevPin{Chip: 1, Line: 27}},  // GPIO1_D3_I2C1_SDA_M4_UART4_RX_M0 - OK
	"32": {Sysfs: 62, Cdev: adaptors.CdevPin{Chip: 1, Line: 30}},  // GPIO1_D6_PWM14_M2_I2C8_SCL_M2 - OK
	"33": {Sysfs: 63, Cdev: adaptors.CdevPin{Chip: 1, Line: 31}},  // GPIO1_D7_PWM15_IR_M3_I2C8_SDA_M2 - OK
	"36": {Sysfs: 131, Cdev: adaptors.CdevPin{Chip: 4, Line: 3}},  // GPIO4_A3_UART0_TX_M2 - PU
	"38": {Sysfs: 132, Cdev: adaptors.CdevPin{Chip: 4, Line: 4}},  // GPIO4_A4_UART0_RX_M2 - OK
	"40": {Sysfs: 133, Cdev: adaptors.CdevPin{Chip: 4, Line: 5}},  // GPIO4_A5_UART3_TX_M2 - OK
	"37": {Sysfs: 134, Cdev: adaptors.CdevPin{Chip: 4, Line: 6}},  // GPIO4_A6_I2C5_SCL_M2_UART3_RX_M2 - OK
	"35": {Sysfs: 135, Cdev: adaptors.CdevPin{Chip: 4, Line: 7}},  // GPIO4_A7_I2C5_SDA_M2 - OK
	"11": {Sysfs: 138, Cdev: adaptors.CdevPin{Chip: 4, Line: 10}}, // GPIO4_B2_CAN1_RX_M1_PWM14_R1 - OK
	"13": {Sysfs: 139, Cdev: adaptors.CdevPin{Chip: 4, Line: 11}}, // GPIO4_B3_CAN1_TX_M1_PWM15_IR_M1 - OK
}

var pwmPinDefinitions = adaptors.PWMPinDefinitions{
	// needs to be enabled by DT-overlay pwm0-m2 (pwm0 = "/pwm@fd8b0000";)
	"27": {Dir: "/sys/devices/platform/fd8b0000.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm1-m2 (pwm0 = "/pwm@fd8b0010";)
	"28": {Dir: "/sys/devices/platform/fd8b0010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm3-m3 (pwm0 = "/pwm@fd8b0030";)
	"12": {Dir: "/sys/devices/platform/fd8b0030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm13-m2  (pwm13 = "/pwm@febf0010";)
	"7": {Dir: "/sys/devices/platform/febf0010.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm14-m1 (pwm14 = "/pwm@febf0020";)
	"11": {Dir: "/sys/devices/platform/febf0020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm14-m2 (pwm14 = "/pwm@febf0020";)
	"32": {Dir: "/sys/devices/platform/febf0020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm15-m1 (pwm15 = "/pwm@febf0030";)
	"13": {Dir: "/sys/devices/platform/febf0030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
	// needs to be enabled by DT-overlay pwm15-m3 (pwm15 = "/pwm@febf0030";)
	"33": {Dir: "/sys/devices/platform/febf0030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3|4|5|6|7]$", Channel: 0},
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
