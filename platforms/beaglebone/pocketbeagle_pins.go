package beaglebone

var pocketBeaglePinMap = map[string]int{
	// P1_01 - VIN
	"P1_02": 87,
	// P1_03 - USB1-V_EN
	"P1_04": 89,
	// P1_05 - USB1-VBUS
	"P1_06": 5,
	// P1_07 - USB1-VIN
	"P1_08": 2,
	// P1_09 - USB1-DN
	"P1_10": 3,
	// P1_11 - USB1-DP
	"P1_12": 4,
	// P1_13 - USB1-ID
	// P1_14 - 3.3V
	// P1_15 - USB1-GND
	// P1_16 - GND
	// P1_16 - AIN-VREF-
	// P1_18 - AIN-VREF+
	// P1_19 - AIO0
	"P1_20": 20,
	// P1_21 - AIO1
	// P1_22 - GND
	// P1_23 - AIO2
	// P1_24 - VOUT-5V
	// P1_25 - AIO3
	"P1_26": 12,
	// P1_27 - AIO4
	"P1_28": 13,
	"P1_29": 117,
	"P1_30": 43,
	"P1_31": 114,
	"P1_32": 42,
	"P1_33": 111,
	"P1_34": 26,
	"P1_35": 88,
	"P1_36": 110,

	"P2_01": 50,
	"P2_02": 59,
	"P2_03": 23,
	"P2_04": 58,
	"P2_05": 30,
	"P2_06": 57,
	"P2_07": 31,
	"P2_08": 60,
	"P2_09": 15,
	"P2_10": 52,
	"P2_11": 14,
	// P2_12 - PWR-BTN
	// P2_13 - VOUT
	// P2_14 - BAT-VIN
	// P2_15 - GND
	// P2_16 - BAT-TEMP
	"P2_17": 65,
	"P2_18": 47,
	"P2_19": 27,
	"P2_20": 64,
	// P2_21 - GND
	"P2_22": 46,
	// P2_23 - 3.3V
	"P2_24": 44,
	"P2_25": 41,
	// P2_26 - NRST
	"P2_27": 40,
	"P2_28": 116,
	"P2_29": 7,
	"P2_30": 113,
	"P2_31": 19,
	"P2_32": 112,
	"P2_33": 45,
	"P2_34": 115,
	"P2_35": 86,
	// P2_36 - AIO7
}

var pocketBeaglePwmPinMap = map[string]pwmPinData{
	"P1_33": {path: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip*", channel: 1},
	"P1_36": {path: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip*", channel: 0},

	"P2_1": {path: "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip*", channel: 0},
	"P2_3": {path: "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip*", channel: 1},
}

var pocketBeagleAnalogPinMap = map[string]string{
	"P1_19": "in_voltage0_raw",
	"P1_21": "in_voltage1_raw",
	"P1_23": "in_voltage2_raw",
	"P1_25": "in_voltage3_raw",
	"P1_27": "in_voltage4_raw",
	"P2_36": "in_voltage7_raw",
}
