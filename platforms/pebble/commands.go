package pebble

func (sd *PebbleDriver) PublishEventC(params map[string]interface{}) {
	sd.PublishEvent(params["name"].(string), params["data"].(string))
}

func (sd *PebbleDriver) SendNotificationC(params map[string]interface{}) {
	sd.SendNotification(params["message"].(string))
}

func (sd *PebbleDriver) PendingMessageC(params map[string]interface{}) interface{} {
	m := make(map[string]string)
	m["result"] = sd.PendingMessage()
	return m
}
