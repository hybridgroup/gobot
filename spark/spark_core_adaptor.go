package spark

import (
	"encoding/json"
	"fmt"
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SparkCoreAdaptor struct {
	gobot.Adaptor
}

func NewSparkCoreAdaptor() *SparkCoreAdaptor {
	return &SparkCoreAdaptor{}
}

func (s *SparkCoreAdaptor) Connect() bool {
	s.Connected = true
	return true
}

func (s *SparkCoreAdaptor) Finalize() bool {
	s.Connected = false
	return true
}

func (s *SparkCoreAdaptor) AnalogRead(pin string) float64 {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.Params["access_token"].(string)},
	}
	url := fmt.Sprintf("%v/analogread", s.deviceUrl())
	resp := s.postToSpark(url, params)
	if resp != nil {
		return resp["return_value"].(float64)
	}
	return 0
}

func (s *SparkCoreAdaptor) PwmWrite(pin string, level byte) {
	s.AnalogWrite(pin, level)
}

func (s *SparkCoreAdaptor) AnalogWrite(pin string, level byte) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, level)},
		"access_token": {s.Params["access_token"].(string)},
	}
	url := fmt.Sprintf("%v/analogwrite", s.deviceUrl())
	s.postToSpark(url, params)
}

func (s *SparkCoreAdaptor) DigitalWrite(pin string, level byte) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, s.pinLevel(level))},
		"access_token": {s.Params["access_token"].(string)},
	}
	url := fmt.Sprintf("%v/digitalwrite", s.deviceUrl())
	s.postToSpark(url, params)
}

func (s *SparkCoreAdaptor) DigitalRead(pin string) int {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.Params["access_token"].(string)},
	}
	url := fmt.Sprintf("%v/digitalread", s.deviceUrl())
	resp := s.postToSpark(url, params)
	if resp != nil {
		return int(resp["return_value"].(float64))
	}
	return -1
}

func (s *SparkCoreAdaptor) deviceUrl() string {
	return fmt.Sprintf("https://api.spark.io/v1/devices/%v", s.Params["device_id"])
}

func (s *SparkCoreAdaptor) pinLevel(level byte) string {
	if level == 1 {
		return "HIGH"
	} else {
		return "LOW"
	}
}

func (s *SparkCoreAdaptor) postToSpark(url string, params url.Values) map[string]interface{} {
	resp, err := http.PostForm(url, params)
	if err != nil {
		fmt.Println(s.Name, "Error writing to spark device", err)
		return nil
	}
	m := make(map[string]interface{})
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(s.Name, "Error reading response body", err)
		return nil
	}
	json.Unmarshal(buf, &m)
	if resp.Status != "200 OK" {
		fmt.Println(s.Name, "Error: ", m["error"])
		return nil
	}
	return m
}
