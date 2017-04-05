package ble

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
)

// GenericAccessDriver represents the Generic Access Service for a BLE Peripheral
type GenericAccessDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewGenericAccessDriver creates a GenericAccessDriver
func NewGenericAccessDriver(a BLEConnector) *GenericAccessDriver {
	n := &GenericAccessDriver{
		name:       gobot.DefaultName("GenericAccess"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the Driver's Connection to the associated Adaptor
func (b *GenericAccessDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver name
func (b *GenericAccessDriver) Name() string { return b.name }

// SetName sets the Driver name
func (b *GenericAccessDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor for this device
func (b *GenericAccessDriver) adaptor() BLEConnector {
	return b.Connection().(BLEConnector)
}

// Start tells driver to get ready to do work
func (b *GenericAccessDriver) Start() (err error) {
	return
}

// Halt stops driver (void)
func (b *GenericAccessDriver) Halt() (err error) { return }

// GetDeviceName returns the device name for the BLE Peripheral
func (b *GenericAccessDriver) GetDeviceName() string {
	c, _ := b.adaptor().ReadCharacteristic("2a00")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetAppearance returns the appearance string for the BLE Peripheral
func (b *GenericAccessDriver) GetAppearance() string {
	c, _ := b.adaptor().ReadCharacteristic("2a01")
	buf := bytes.NewBuffer(c)

	var val uint16
	binary.Read(buf, binary.LittleEndian, &val)
	return appearances[val]
}

var appearances = map[uint16]string{
	0:    "Unknown",
	1024: "Generic Glucose Meter",
	1088: "Generic: Running Walking Sensor",
	1089: "Running Walking Sensor: In-Shoe",
	1090: "Running Walking Sensor: On-Shoe",
	1091: "Running Walking Sensor: On-Hip",
	1152: "Generic: Cycling",
	1153: "Cycling: Cycling Computer",
	1154: "Cycling: Speed Sensor",
	1155: "Cycling: Cadence Sensor",
	1156: "Cycling: Power Sensor",
	1157: "Cycling: Speed and Cadence Sensor",
	128:  "Generic Computer",
	192:  "Generic Watch",
	193:  "Watch: Sports Watch",
	256:  "Generic Clock",
	3136: "Generic: Pulse Oximeter",
	3137: "Fingertip Pulse",
	3138: "Wrist Worn",
	320:  "Generic Display",
	3200: "Generic: Weight Scale",
	384:  "Generic Remote Control",
	448:  "Generic Eye-glasses",
	512:  "Generic Tag",
	5184: "Generic: Outdoor Sports Activity",
	5185: "Location Display Device",
	5186: "Location and Navigation Display Device",
	5187: "Location Pod",
	5188: "Location and Navigation Pod",
	576:  "Generic Keyring",
	64:   "Generic Phone",
	640:  "Generic Media Player",
	704:  "Generic Barcode Scanner",
	768:  "Generic Thermometer",
	769:  "Thermometer: Ear",
	832:  "Generic Heart rate Sensor",
	833:  "Heart Rate Sensor: Heart Rate Belt",
	896:  "Generic Blood Pressure",
	897:  "Blood Pressure: Arm Blood",
	898:  "Blood Pressure: Wrist Blood",
	960:  "Human Interface Device (HID)",
	961:  "Keyboard",
	962:  "Mouse",
	963:  "Joystick",
	964:  "Gamepad",
	965:  "Digitizer Tablet",
	966:  "Card Reader",
	967:  "Digital Pen",
	968:  "Barcode Scanner",
}
