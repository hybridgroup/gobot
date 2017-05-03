package beaglebone

type pwmPinData struct {
	channel int
	path    string
}

var pins = map[string]int{
	// P8_1 - P8_2 GND
	// P8_3 - P8_6 EMCC
	"P8_7":  66,
	"P8_8":  67,
	"P8_9":  69,
	"P8_10": 68,
	"P8_11": 45,
	"P8_12": 44,
	"P8_13": 23,
	"P8_14": 26,
	"P8_15": 47,
	"P8_16": 46,
	"P8_17": 27,
	"P8_18": 65,
	"P8_19": 22,
	// P8_20 - P8_25 EMCC
	"P8_26": 61,
	// P8_27 - P8_46 HDMI

	// P9_1 - P9_2 GND
	// P9_3 - P9_4 3V3
	// P9_5 - P9_6 5V
	// P9_7 - P9_8 5V SYS
	// P9_9 PWR_BUT
	// P9_10 SYS_RESET
	"P9_11": 30,
	"P9_12": 60,
	"P9_13": 31,
	"P9_14": 50,
	"P9_15": 48,
	"P9_16": 51,
	"P9_17": 5,
	"P9_18": 4,
	// P9_19 I2C2 SCL
	// P9_20 I2C2 SDA
	"P9_21": 3,
	"P9_22": 2,
	"P9_23": 49,
	"P9_24": 15,
	"P9_25": 117,
	"P9_26": 14,
	"P9_27": 115,
	"P9_28": 113,
	"P9_29": 111,
	"P9_30": 112,
	"P9_31": 110,
}

var pwmPins = map[string]pwmPinData{
	"P8_13": {path: "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4", channel: 1},
	"P8_19": {path: "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4", channel: 0},

	"P9_14": {path: "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2", channel: 0},
	"P9_16": {path: "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2", channel: 1},
	"P9_21": {path: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0", channel: 1},
	"P9_22": {path: "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0", channel: 0},
	//"P9_42": {path: "", channel: 0}, TODO: implement this pwm
}

var analogPins = map[string]string{
	"P9_39": "in_voltage0_raw",
	"P9_40": "in_voltage1_raw",
	"P9_37": "in_voltage2_raw",
	"P9_38": "in_voltage3_raw",
	"P9_33": "in_voltage4_raw",
	"P9_36": "in_voltage5_raw",
	"P9_35": "in_voltage6_raw",
}
