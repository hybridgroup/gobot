package gobotI2C

// blinkm
func (self *BlinkM) FirmwareVersionC(params map[string]interface{}) string {
	return self.FirmwareVersion()
}
func (self *BlinkM) ColorC(params map[string]interface{}) []byte {
	return self.Color()
}
func (self *BlinkM) RgbC(params map[string]interface{}) {
	r := byte(params["r"].(float64))
	g := byte(params["g"].(float64))
	b := byte(params["b"].(float64))
	self.Rgb(r, g, b)
}
func (self *BlinkM) FadeC(params map[string]interface{}) {
	r := byte(params["r"].(float64))
	g := byte(params["g"].(float64))
	b := byte(params["b"].(float64))
	self.Fade(r, g, b)
}
