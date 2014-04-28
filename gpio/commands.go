package gpio

import "strconv"

// Led
func (l *LedDriver) ToggleC(params map[string]interface{}) {
	l.Toggle()
}
func (l *LedDriver) OnC(params map[string]interface{}) {
	l.On()
}
func (l *LedDriver) OffC(params map[string]interface{}) {
	l.Off()
}
func (l *LedDriver) BrightnessC(params map[string]interface{}) {
	level := byte(params["level"].(float64))
	l.Brightness(level)
}

// Servo
func (l *ServoDriver) MoveC(params map[string]interface{}) {
	angle := byte(params["angle"].(float64))
	l.Move(angle)
}
func (l *ServoDriver) MinC(params map[string]interface{}) {
	l.Min()
}
func (l *ServoDriver) CenterC(params map[string]interface{}) {
	l.Center()
}
func (l *ServoDriver) MaxC(params map[string]interface{}) {
	l.Max()
}

// Direct Pin
func (d *DirectPinDriver) DigitalReadC(params map[string]interface{}) int {
	return d.DigitalRead()
}
func (d *DirectPinDriver) DigitalWriteC(params map[string]interface{}) {
	level, _ := strconv.Atoi(params["level"].(string))
	d.DigitalWrite(byte(level))
}
func (d *DirectPinDriver) AnalogReadC(params map[string]interface{}) int {
	return d.AnalogRead()
}
func (d *DirectPinDriver) AnalogWriteC(params map[string]interface{}) {
	level, _ := strconv.Atoi(params["level"].(string))
	d.AnalogWrite(byte(level))
}
func (d *DirectPinDriver) PwmWriteC(params map[string]interface{}) {
	level, _ := strconv.Atoi(params["level"].(string))
	d.PwmWrite(byte(level))
}
func (d *DirectPinDriver) ServoWriteC(params map[string]interface{}) {
	level, _ := strconv.Atoi(params["level"].(string))
	d.ServoWrite(byte(level))
}

// Analog Sensor
func (d *AnalogSensorDriver) ReadC(params map[string]interface{}) int {
	return d.Read()
}
