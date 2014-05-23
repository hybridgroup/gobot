package sphero

func (s *SpheroDriver) SetRGBC(params map[string]interface{}) {
	r := uint8(params["r"].(float64))
	g := uint8(params["g"].(float64))
	b := uint8(params["b"].(float64))
	s.SetRGB(r, g, b)
}

func (s *SpheroDriver) RollC(params map[string]interface{}) {
	speed := uint8(params["speed"].(float64))
	heading := uint16(params["heading"].(float64))
	s.Roll(speed, heading)
}

func (s *SpheroDriver) StopC() {
	s.Stop()
}

func (s *SpheroDriver) GetRGBC() {
}

func (s *SpheroDriver) SetBackLEDC(params map[string]interface{}) {
	level := uint8(params["level"].(float64))
	s.SetBackLED(level)
}

func (s *SpheroDriver) SetHeadingC(params map[string]interface{}) {
	heading := uint16(params["heading"].(float64))
	s.SetHeading(heading)
}
func (s *SpheroDriver) SetStabilizationC(params map[string]interface{}) {
	on := params["heading"].(bool)
	s.SetStabilization(on)
}
