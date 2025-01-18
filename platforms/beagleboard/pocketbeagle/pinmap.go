package pocketbeagle

import "gobot.io/x/gobot/v2/platforms/adaptors"

// tested: am335x-debian-11.7-iot-armhf-2023-09-02-4gb, OK: work as input and output, IN: work only as input
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	// gpiochip0 - 32 lines:
	// line   0:         "NC"       unused   input  active-high
	// line   1:         "NC"       unused   input  active-high
	"P1_08": {Sysfs: 2, Cdev: adaptors.CdevPin{Chip: 0, Line: 2}}, // P1.08_SPI0_CLK - ?
	"P1_10": {Sysfs: 3, Cdev: adaptors.CdevPin{Chip: 0, Line: 3}}, // P1.10_SPI0_MISO - ?
	"P1_12": {Sysfs: 4, Cdev: adaptors.CdevPin{Chip: 0, Line: 4}}, // P1.12_SPI0_MOSI - ?
	"P1_06": {Sysfs: 5, Cdev: adaptors.CdevPin{Chip: 0, Line: 5}}, // P1.06_SPI0_CS - ?
	// line   6:  "[MMC0_CD]"         "cd"   input   active-low [used]
	"P2_29": {Sysfs: 7, Cdev: adaptors.CdevPin{Chip: 0, Line: 7}}, // P2.29_SPI1_CLK - ?
	// line   8: "[SYSBOOT 12]" unused input active-high
	// line   9: "[SYSBOOT 13]" unused input active-high
	// line  10: "[SYSBOOT 14]" unused input active-high
	// line  11: "[SYSBOOT 15]" unused input active-high
	"P1_26": {Sysfs: 12, Cdev: adaptors.CdevPin{Chip: 0, Line: 12}}, // P1.26_I2C2_SDA - ?
	"P1_28": {Sysfs: 13, Cdev: adaptors.CdevPin{Chip: 0, Line: 13}}, // P1.28_I2C2_SCL - ?
	"P2_11": {Sysfs: 14, Cdev: adaptors.CdevPin{Chip: 0, Line: 14}}, // P2.11_I2C1_SDA - ?
	"P2_09": {Sysfs: 15, Cdev: adaptors.CdevPin{Chip: 0, Line: 15}}, // P2.09_I2C1_SCL - ?
	// line  16:         "NC"       unused   input  active-high
	// line  17:         "NC"       unused   input  active-high
	// line  18:         "NC"       unused   input  active-high
	"P2_31": {Sysfs: 19, Cdev: adaptors.CdevPin{Chip: 0, Line: 19}}, // P2.31_SPI1_CS - ?
	"P1_20": {Sysfs: 20, Cdev: adaptors.CdevPin{Chip: 0, Line: 20}}, // P1.20_PRU0.16 - ?
	// line  21:         "NC"       unused   input  active-high
	// line  22:         "NC"       unused   input  active-high
	"P2_03": {Sysfs: 23, Cdev: adaptors.CdevPin{Chip: 0, Line: 23}}, // P2.03 - ?
	// line  24:         "NC"       unused   input  active-high
	// line  25:         "NC"       unused   input  active-high
	"P1_34": {Sysfs: 26, Cdev: adaptors.CdevPin{Chip: 0, Line: 26}}, // P1.34 - ?
	"P2_19": {Sysfs: 27, Cdev: adaptors.CdevPin{Chip: 0, Line: 27}}, // P2.19 - ?
	// line  28:         "NC"       unused   input  active-high
	// line  29:         "NC"       unused   input  active-high
	"P2_05": {Sysfs: 30, Cdev: adaptors.CdevPin{Chip: 0, Line: 30}}, // P2.05_UART4_RX - ?
	"P2_07": {Sysfs: 31, Cdev: adaptors.CdevPin{Chip: 0, Line: 31}}, // P2.07_UART4_TX - ?
	// gpiochip1 - 32 lines:
	// line   0:         "NC"       unused   input  active-high
	// ...
	// line   7:         "NC"       unused   input  active-high
	"P2_27": {Sysfs: 40, Cdev: adaptors.CdevPin{Chip: 1, Line: 8}},  // P2.27_SPI1_MISO- ?
	"P2_25": {Sysfs: 41, Cdev: adaptors.CdevPin{Chip: 1, Line: 9}},  // P2.25_SPI1_MOSI - ?
	"P1_32": {Sysfs: 42, Cdev: adaptors.CdevPin{Chip: 1, Line: 10}}, // P1.32_UART0_RX - ?
	"P1_30": {Sysfs: 43, Cdev: adaptors.CdevPin{Chip: 1, Line: 11}}, // P1.30_UART0_TX - ?
	"P2_24": {Sysfs: 44, Cdev: adaptors.CdevPin{Chip: 1, Line: 12}}, // P2.24 - ?
	"P2_33": {Sysfs: 45, Cdev: adaptors.CdevPin{Chip: 1, Line: 13}}, // P2.33 - ?
	"P2_22": {Sysfs: 46, Cdev: adaptors.CdevPin{Chip: 1, Line: 14}}, // P2.22 - ?
	"P2_18": {Sysfs: 47, Cdev: adaptors.CdevPin{Chip: 1, Line: 15}}, // P2.18 - ?
	// line  16:         "NC"       unused   input  active-high
	// line  17:         "NC"       unused   input  active-high
	"P2_01": {Sysfs: 50, Cdev: adaptors.CdevPin{Chip: 1, Line: 18}}, // P2.01_PWM1A - ?
	// line  19:         "NC"       unused   input  active-high
	"P2_10": {Sysfs: 52, Cdev: adaptors.CdevPin{Chip: 1, Line: 20}}, // P2.10- ?
	"usr0":  {Sysfs: -1, Cdev: adaptors.CdevPin{Chip: 1, Line: 21}}, // USR LED 0 (beaglebone:green:usr0) - ?
	"usr1":  {Sysfs: -1, Cdev: adaptors.CdevPin{Chip: 1, Line: 22}}, // USR LED 1 (beaglebone:green:usr1) - ?
	"usr2":  {Sysfs: -1, Cdev: adaptors.CdevPin{Chip: 1, Line: 23}}, // USR LED 2 (beaglebone:green:usr2) - ?
	"usr3":  {Sysfs: -1, Cdev: adaptors.CdevPin{Chip: 1, Line: 24}}, // USR LED 3 (beaglebone:green:usr3) - ?
	"P2_06": {Sysfs: 57, Cdev: adaptors.CdevPin{Chip: 1, Line: 25}}, // P2.06 - ?
	"P2_04": {Sysfs: 58, Cdev: adaptors.CdevPin{Chip: 1, Line: 26}}, // P2.04 - ?
	"P2_02": {Sysfs: 59, Cdev: adaptors.CdevPin{Chip: 1, Line: 27}}, // P2.02 - ?
	"P2_08": {Sysfs: 60, Cdev: adaptors.CdevPin{Chip: 1, Line: 28}}, // P2.08 - ?
	// line  29:         "NC"       unused   input  active-high
	// line  30:         "NC"       unused   input  active-high
	// line  31:         "NC"       unused   input  active-high
	// gpiochip2 - 32 lines:
	"P2_20": {Sysfs: 64, Cdev: adaptors.CdevPin{Chip: 2, Line: 0}}, // P2.20 - ?
	"P2_17": {Sysfs: 65, Cdev: adaptors.CdevPin{Chip: 2, Line: 1}}, // P2.17 - ?
	// line   2:         "NC"       unused   input  active-high
	// line   3:         "NC"       unused   input  active-high
	// line   4:         "NC"       unused   input  active-high
	// line   5: "[EEPROM_WP]" unused input active-high
	// line   6: "[SYSBOOT 0]" unused input active-high
	// line   7: "[SYSBOOT 1]" unused input active-high
	// line   8: "[SYSBOOT 2]" unused input active-high
	// line   9: "[SYSBOOT 3]" unused input active-high
	// line  10: "[SYSBOOT 4]" unused input active-high
	// line  11: "[SYSBOOT 5]" unused input active-high
	// line  12: "[SYSBOOT 6]" unused input active-high
	// line  13: "[SYSBOOT 7]" unused input active-high
	// line  14: "[SYSBOOT 8]" unused input active-high
	// line  15: "[SYSBOOT 9]" unused input active-high
	// line  16: "[SYSBOOT 10]" unused input active-high
	// line  17: "[SYSBOOT 11]" unused input active-high
	// line  18:         "NC"       unused   input  active-high
	// ...
	// line  21:         "NC"       unused   input  active-high
	"P2_35": {Sysfs: 86, Cdev: adaptors.CdevPin{Chip: 2, Line: 22}}, // P2.35_AIN5 - ?
	"P1_02": {Sysfs: 87, Cdev: adaptors.CdevPin{Chip: 2, Line: 23}}, // P1.02_AIN6 - ?
	"P1_35": {Sysfs: 88, Cdev: adaptors.CdevPin{Chip: 2, Line: 24}}, // P1.35_PRU1.10 - ?
	"P1_04": {Sysfs: 89, Cdev: adaptors.CdevPin{Chip: 2, Line: 25}}, // P1.04_PRU1.11 - ?
	// line  26: "[MMC0_DAT3]" unused input active-high
	// line  27: "[MMC0_DAT2]" unused input active-high
	// line  28: "[MMC0_DAT1]" unused input active-high
	// line  29: "[MMC0_DAT0]" unused input active-high
	// line  30: "[MMC0_CLK]"       unused   input  active-high
	// line  31: "[MMC0_CMD]"       unused   input  active-high
	// gpiochip3 - 32 lines:
	// line   0:         "NC"       unused   input  active-high
	// ...
	// line   4:         "NC"       unused   input  active-high
	// line  13: "P1.03 [USB1]" unused input active-high
	"P1_36": {Sysfs: 110, Cdev: adaptors.CdevPin{Chip: 3, Line: 14}}, // P1.36_PWM0A - ?
	"P1_33": {Sysfs: 111, Cdev: adaptors.CdevPin{Chip: 3, Line: 15}}, // P1.33_PRU0.1 - ?
	"P2_32": {Sysfs: 112, Cdev: adaptors.CdevPin{Chip: 3, Line: 16}}, // P2.32_PRU0.2 - ?
	"P2_30": {Sysfs: 113, Cdev: adaptors.CdevPin{Chip: 3, Line: 17}}, // P2.30_PRU0.3 - ?
	"P1_31": {Sysfs: 114, Cdev: adaptors.CdevPin{Chip: 3, Line: 18}}, // P1.31_PRU0.4 - ?
	"P2_34": {Sysfs: 115, Cdev: adaptors.CdevPin{Chip: 3, Line: 19}}, // P2.34_PRU0.5 - ?
	"P2_28": {Sysfs: 116, Cdev: adaptors.CdevPin{Chip: 3, Line: 20}}, // P2.28_PRU0.6- ?
	"P1_29": {Sysfs: 117, Cdev: adaptors.CdevPin{Chip: 3, Line: 21}}, // P1.29_PRU0.7 - ?
	// line  22:         "NC"       unused   input  active-high
	// ...
	// line  31:         "NC"       unused   input  active-high

	// P1_01 - VIN; P1_03 - USB1-V_EN; P1_05 - USB1-VBUS; P1_07 - USB1-VIN; P1_09 - USB1-DN; P1_11 - USB1-DP;
	// P1_13 - USB1-ID; P1_14 - 3.3V; P1_15 - USB1-GND; P1_16 - GND; P1_16 - AIN-VREF-; P1_18 - AIN-VREF+; P1_19 - AIO0
	// P1_21 - AIO1; P1_22 - GND; P1_23 - AIO2; P1_24 - VOUT-5V; P1_25 - AIO3; P1_27 - AIO4; P2_12 - PWR-BTN; P2_13 - VOUT
	// P2_14 - BAT-VIN; P2_15 - GND; P2_16 - BAT-TEMP; P2_21 - GND; P2_23 - 3.3V; P2_26 - NRST; P2_36 - AIO7
}

var pwmPinMap = adaptors.PWMPinDefinitions{
	"P1_33": {Dir: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/", DirRegexp: "pwmchip[0-9]+$", Channel: 1},
	"P1_36": {Dir: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/", DirRegexp: "pwmchip[0-9]+$", Channel: 0},

	"P2_1": {Dir: "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/", DirRegexp: "pwmchip[0-9]+$", Channel: 0},
	"P2_3": {Dir: "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/", DirRegexp: "pwmchip[0-9]+$", Channel: 1},
}

var analogPinMap = adaptors.AnalogPinDefinitions{
	"P1_19": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage0_raw", W: false, ReadBufLen: 1024},
	"P1_21": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage1_raw", W: false, ReadBufLen: 1024},
	"P1_23": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage2_raw", W: false, ReadBufLen: 1024},
	"P1_25": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage3_raw", W: false, ReadBufLen: 1024},
	"P1_27": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage4_raw", W: false, ReadBufLen: 1024},
	"P2_36": {Path: "/sys/bus/iio/devices/iio:device0/in_voltage7_raw", W: false, ReadBufLen: 1024},
}
