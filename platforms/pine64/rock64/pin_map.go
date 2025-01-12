package rock64

import "gobot.io/x/gobot/v2/platforms/adaptors"

// notes for character device
// sysfs: Chip*32 + (A=0, B=8, C=16) + Nr
// tested with cdev on a ROCK64 V2.0 board: armbian Linux, OK: works, ?: unknown, NOK: not working
// IN: works only as input, PU: if used as input, external pullup resistor needed
// PM: pins are shared with the PMIC i2c communication (address 0x13), GPIO seems to work but problems can occur
//
// pin suffix:
// "P5": those pins are located on the second pin header "P5+BUS"
// "SD": those pins should be used with caution and only if eMMC is used and no SD card is inserted
// "M2": those pins are used by the SPI FLASH 128M memory chip, but not blocked and maybe work, if unsure don't use it
// "SDA/SCL": see PM above, should only used for i2c-1
var gpioPinDefinitions = adaptors.DigitalPinDefinitions{
	"13":     {Sysfs: 0, Cdev: adaptors.CdevPin{Chip: 0, Line: 0}},   // GPIO0_A0 - OK
	"P5_13":  {Sysfs: 27, Cdev: adaptors.CdevPin{Chip: 0, Line: 27}}, // GPIO0_D3_SPDIF_TX_M0 - OK
	"33_SD":  {Sysfs: 32, Cdev: adaptors.CdevPin{Chip: 1, Line: 0}},  // GPIO1_A0_SDMMC0_D0 - ?
	"35_SD":  {Sysfs: 33, Cdev: adaptors.CdevPin{Chip: 1, Line: 1}},  // GPIO1_A1_SDMMC0_D1 - ?
	"37_SD":  {Sysfs: 34, Cdev: adaptors.CdevPin{Chip: 1, Line: 2}},  // GPIO1_A2_SDMMC0_D2_JTAG_TCK - ?
	"40_SD":  {Sysfs: 35, Cdev: adaptors.CdevPin{Chip: 1, Line: 3}},  // GPIO1_A3_SDMMC0_D3_JTAG_TMS - ?
	"38_SD":  {Sysfs: 36, Cdev: adaptors.CdevPin{Chip: 1, Line: 4}},  // GPIO1_A4_SDMMC0_CMD - ?
	"36_SD":  {Sysfs: 37, Cdev: adaptors.CdevPin{Chip: 1, Line: 5}},  // GPIO1_A5_SDMMC0_DET - ?
	"32_SD":  {Sysfs: 38, Cdev: adaptors.CdevPin{Chip: 1, Line: 6}},  // GPIO1_A6_SDMMC0_CLK - ?
	"7":      {Sysfs: 60, Cdev: adaptors.CdevPin{Chip: 1, Line: 28}}, // GPIO1_D4_CLK32KOUT_M1 - OK (PU)
	"8":      {Sysfs: 64, Cdev: adaptors.CdevPin{Chip: 2, Line: 0}},  // GPIO2_A0_UART2_TX_M1 - OK (PU)
	"10":     {Sysfs: 65, Cdev: adaptors.CdevPin{Chip: 2, Line: 1}},  // GPIO2_A1_UART2_RX_M1 - OK
	"12":     {Sysfs: 67, Cdev: adaptors.CdevPin{Chip: 2, Line: 3}},  // GPIO2_A3 - IN
	"27_SDA": {Sysfs: 68, Cdev: adaptors.CdevPin{Chip: 2, Line: 4}},  // GPIO2_A4_I2C1_SDA - PM
	"28_SCL": {Sysfs: 69, Cdev: adaptors.CdevPin{Chip: 2, Line: 5}},  // GPIO2_A5_I2C1_SCL - PM
	"26":     {Sysfs: 76, Cdev: adaptors.CdevPin{Chip: 2, Line: 12}}, // GPIO2_B4_SPI_CSN1_M0 - OK
	"P5_10":  {Sysfs: 79, Cdev: adaptors.CdevPin{Chip: 2, Line: 15}}, // GPIO2_B7_I2S1_MCLK - OK (PU)
	"P5_9":   {Sysfs: 80, Cdev: adaptors.CdevPin{Chip: 2, Line: 16}}, // GPIO2_C0_I2S1_LRCKRX - OK
	"P5_3":   {Sysfs: 81, Cdev: adaptors.CdevPin{Chip: 2, Line: 17}}, // GPIO2_C1_I2S1_LRCKTX - OK
	"P5_4":   {Sysfs: 82, Cdev: adaptors.CdevPin{Chip: 2, Line: 18}}, // GPIO2_C2_I2S1_SCLK - OK (PU)
	"P5_6":   {Sysfs: 83, Cdev: adaptors.CdevPin{Chip: 2, Line: 19}}, // GPIO2_C3_I2S1_SDI - OK
	"P5_12":  {Sysfs: 84, Cdev: adaptors.CdevPin{Chip: 2, Line: 20}}, // GPIO2_C4_I2S1_SDIO1 - OK
	"P5_11":  {Sysfs: 85, Cdev: adaptors.CdevPin{Chip: 2, Line: 21}}, // GPIO2_C5_I2S1_SDIO2 - OK
	"P5_14":  {Sysfs: 86, Cdev: adaptors.CdevPin{Chip: 2, Line: 22}}, // GPIO2_C6_I2S1_SDIO3 - OK
	"P5_5":   {Sysfs: 87, Cdev: adaptors.CdevPin{Chip: 2, Line: 23}}, // GPIO2_C7_I2S1_SDO - OK
	"5":      {Sysfs: 88, Cdev: adaptors.CdevPin{Chip: 2, Line: 24}}, // GPIO2_D0_I2C0_SCL_EthernetLink (P5_22) - OK
	"P5_22":  {Sysfs: 88, Cdev: adaptors.CdevPin{Chip: 2, Line: 24}}, // GPIO2_D0_I2C0_SCL_EthernetLink (P5_22) - OK
	"3":      {Sysfs: 89, Cdev: adaptors.CdevPin{Chip: 2, Line: 25}}, // GPIO2_D1_I2C0_SDA_EthernetSpeed (P5_21) - OK
	"P5_21":  {Sysfs: 89, Cdev: adaptors.CdevPin{Chip: 2, Line: 25}}, // GPIO2_D1_I2C0_SDA_EthernetSpeed (P5_21) - OK
	"23_M2":  {Sysfs: 96, Cdev: adaptors.CdevPin{Chip: 3, Line: 0}},  // GPIO3_A0_SPI_CLK_M2 - NOK
	"19_M2":  {Sysfs: 97, Cdev: adaptors.CdevPin{Chip: 3, Line: 1}},  // GPIO3_A1_SPI_TXD_M2 - OK
	"21_M2":  {Sysfs: 98, Cdev: adaptors.CdevPin{Chip: 3, Line: 2}},  // GPIO3_A2_SPI_RXD_M2 - OK
	"15":     {Sysfs: 100, Cdev: adaptors.CdevPin{Chip: 3, Line: 4}}, // GPIO3_A4 - OK
	"16":     {Sysfs: 101, Cdev: adaptors.CdevPin{Chip: 3, Line: 5}}, // GPIO3_A5 - OK
	"18":     {Sysfs: 102, Cdev: adaptors.CdevPin{Chip: 3, Line: 6}}, // GPIO3_A6 - OK
	"22":     {Sysfs: 103, Cdev: adaptors.CdevPin{Chip: 3, Line: 7}}, // GPIO3_A7 - OK
	"24_M2":  {Sysfs: 104, Cdev: adaptors.CdevPin{Chip: 3, Line: 8}}, // GPIO3_B0_SPI_CSN0_M2 - NOK
}

var analogPinDefinitions = adaptors.AnalogPinDefinitions{
	// +/-273.200 °C need >=7 characters to read: +/-273200 millidegree Celsius
	"thermal_zone0": {Path: "/sys/class/thermal/thermal_zone0/temp", W: false, ReadBufLen: 7},
}
