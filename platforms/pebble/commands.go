package pebble

func (sd *PebbleDriver) PublishEventC(params map[string]interface{}) {
	sd.PublishEvent(params["name"].(string), params["data"].(string))
}
