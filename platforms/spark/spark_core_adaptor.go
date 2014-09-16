package spark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hybridgroup/gobot"
)

type SparkCoreAdaptor struct {
	gobot.Adaptor
	DeviceID    string
	AccessToken string
	APIServer   string
}

func NewSparkCoreAdaptor(name string, deviceID string, accessToken string) *SparkCoreAdaptor {
	return &SparkCoreAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"SparkCoreAdaptor",
		),
		DeviceID:    deviceID,
		AccessToken: accessToken,
		APIServer:   "https://api.spark.io",
	}
}

func (s *SparkCoreAdaptor) Connect() bool {
	s.SetConnected(true)
	return true
}

func (s *SparkCoreAdaptor) Finalize() bool {
	s.SetConnected(false)
	return true
}

func (s *SparkCoreAdaptor) AnalogRead(pin string) int {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}

	url := fmt.Sprintf("%v/analogread", s.deviceURL())

	resp, err := s.postToSpark(url, params)
	if err == nil {
		return int(resp["return_value"].(float64))
	}

	return 0
}

func (s *SparkCoreAdaptor) PwmWrite(pin string, level byte) {
	s.AnalogWrite(pin, level)
}

func (s *SparkCoreAdaptor) AnalogWrite(pin string, level byte) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, level)},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/analogwrite", s.deviceURL())
	s.postToSpark(url, params)
}

func (s *SparkCoreAdaptor) DigitalWrite(pin string, level byte) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, s.pinLevel(level))},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalwrite", s.deviceURL())
	s.postToSpark(url, params)
}

func (s *SparkCoreAdaptor) DigitalRead(pin string) int {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalread", s.deviceURL())
	resp, err := s.postToSpark(url, params)
	if err == nil {
		return int(resp["return_value"].(float64))
	}
	return -1
}

func (s *SparkCoreAdaptor) setAPIServer(server string) {
	s.APIServer = server
}

func (s *SparkCoreAdaptor) deviceURL() string {
	if len(s.APIServer) <= 0 {
		s.setAPIServer("https://api.spark.io")
	}
	return fmt.Sprintf("%v/v1/devices/%v", s.APIServer, s.DeviceID)
}

func (s *SparkCoreAdaptor) pinLevel(level byte) string {
	if level == 1 {
		return "HIGH"
	}
	return "LOW"
}

func (s *SparkCoreAdaptor) postToSpark(url string, params url.Values) (m map[string]interface{}, err error) {
	resp, err := http.PostForm(url, params)
	if err != nil {
		fmt.Println(s.Name, "Error writing to spark device", err)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(s.Name, "Error reading response body", err)
		return
	}

	json.Unmarshal(buf, &m)

	if resp.Status != "200 OK" {
		fmt.Println(s.Name, "Error: ", m["error"])
		err = fmt.Errorf("%q was not found", url)
		return
	}

	return
}
