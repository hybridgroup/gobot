package i2c

// blinkm
func (b *BlinkMDriver) FirmwareVersionC(params map[string]interface{}) string {
	return b.FirmwareVersion()
}
func (b *BlinkMDriver) ColorC(params map[string]interface{}) []byte {
	return b.Color()
}
func (b *BlinkMDriver) RgbC(params map[string]interface{}) {
	red := byte(params["red"].(float64))
	green := byte(params["green"].(float64))
	blue := byte(params["blue"].(float64))
	b.Rgb(red, green, blue)
}
func (b *BlinkMDriver) FadeC(params map[string]interface{}) {
	red := byte(params["red"].(float64))
	green := byte(params["green"].(float64))
	blue := byte(params["blue"].(float64))
	b.Fade(red, green, blue)
}
