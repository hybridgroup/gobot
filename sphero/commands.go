package gobotSphero

func (sd *SpheroDriver) SetRGBC(params map[string]interface{}) {
	r := uint8(params["r"].(float64))
	g := uint8(params["g"].(float64))
	b := uint8(params["b"].(float64))
	sd.SetRGB(r, g, b)
}

func (sd *SpheroDriver) RollC(params map[string]interface{}) {
	speed := uint8(params["speed"].(float64))
	heading := uint16(params["heading"].(float64))
	sd.Roll(speed, heading)
}

func (sd *SpheroDriver) StopC() {
	sd.Stop()
}

func (sd *SpheroDriver) GetRGBC() {
}

func (sd *SpheroDriver) SetBackLEDC(params map[string]interface{}) {
	level := uint8(params["level"].(float64))
	sd.SetBackLED(level)
}

func (sd *SpheroDriver) SetHeadingC(params map[string]interface{}) {
	heading := uint16(params["heading"].(float64))
	sd.SetHeading(heading)
}
func (sd *SpheroDriver) SetStabilizationC(params map[string]interface{}) {
	on := params["heading"].(bool)
	sd.SetStabilization(on)
}
