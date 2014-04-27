package sphero

func (s *Sphero) SetRGBC(params map[string]interface{}) {
	r := uint8(params["r"].(float64))
	g := uint8(params["g"].(float64))
	b := uint8(params["b"].(float64))
	s.SetRGB(r, g, b)
}

func (s *Sphero) RollC(params map[string]interface{}) {
	speed := uint8(params["speed"].(float64))
	heading := uint16(params["heading"].(float64))
	s.Roll(speed, heading)
}

func (s *Sphero) StopC() {
	s.Stop()
}

func (s *Sphero) GetRGBC() {
}

func (s *Sphero) SetBackLEDC(params map[string]interface{}) {
	level := uint8(params["level"].(float64))
	s.SetBackLED(level)
}

func (s *Sphero) SetHeadingC(params map[string]interface{}) {
	heading := uint16(params["heading"].(float64))
	s.SetHeading(heading)
}
func (s *Sphero) SetStabilizationC(params map[string]interface{}) {
	on := params["heading"].(bool)
	s.SetStabilization(on)
}
