package spark

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*SparkCoreAdaptor)(nil)

type SparkCoreAdaptor struct {
	gobot.Adaptor
	DeviceID    string
	AccessToken string
	APIServer   string
}

// NewSparkCoreAdaptor creates new spark core adaptor with deviceId and accessToken
// using api.spark.io server as default
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

// Connect returns true if connection to spark core is succesfull
func (s *SparkCoreAdaptor) Connect() error {
	s.SetConnected(true)
	return nil
}

// Finalize returns true if connection to spark core is finalized successfully
func (s *SparkCoreAdaptor) Finalize() error {
	s.SetConnected(false)
	return nil
}

// AnalogRead reads analog ping value using spark cloud api
func (s *SparkCoreAdaptor) AnalogRead(pin string) (val int, err error) {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}

	url := fmt.Sprintf("%v/analogread", s.deviceURL())

	resp, err := s.postToSpark(url, params)
	if err == nil {
		val = int(resp["return_value"].(float64))
		return
	}

	return 0, err
}

// PwmWrite writes in pin using analog write api
func (s *SparkCoreAdaptor) PwmWrite(pin string, level byte) (err error) {
	return s.AnalogWrite(pin, level)
}

// AnalogWrite writes analog pin with specified level using spark cloud api
func (s *SparkCoreAdaptor) AnalogWrite(pin string, level byte) (err error) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, level)},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/analogwrite", s.deviceURL())
	_, err = s.postToSpark(url, params)
	return
}

// DigitalWrite writes to a digital pin using spark cloud api
func (s *SparkCoreAdaptor) DigitalWrite(pin string, level byte) (err error) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, s.pinLevel(level))},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalwrite", s.deviceURL())
	_, err = s.postToSpark(url, params)
	return err
}

// DigitalRead reads from digital pin using spark cloud api
func (s *SparkCoreAdaptor) DigitalRead(pin string) (val int, err error) {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalread", s.deviceURL())
	resp, err := s.postToSpark(url, params)
	if err == nil {
		val = int(resp["return_value"].(float64))
		return
	}
	return -1, err
}

// setAPIServer sets spark cloud api server, this can be used to change from default api.spark.io
func (s *SparkCoreAdaptor) setAPIServer(server string) {
	s.APIServer = server
}

// deviceURL constructs device url to make requests from spark cloud api
func (s *SparkCoreAdaptor) deviceURL() string {
	if len(s.APIServer) <= 0 {
		s.setAPIServer("https://api.spark.io")
	}
	return fmt.Sprintf("%v/v1/devices/%v", s.APIServer, s.DeviceID)
}

// pinLevel converts byte level to string expected in api
func (s *SparkCoreAdaptor) pinLevel(level byte) string {
	if level == 1 {
		return "HIGH"
	}
	return "LOW"
}

// postToSpark makes POST request to spark cloud server, return err != nil if there is
// any issue with the request.
func (s *SparkCoreAdaptor) postToSpark(url string, params url.Values) (m map[string]interface{}, err error) {
	resp, err := http.PostForm(url, params)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	json.Unmarshal(buf, &m)

	if resp.Status != "200 OK" {
		if _, ok := m["error"]; ok {
			err = errors.New(m["error"].(string))
		} else {
			err = errors.New(fmt.Sprintf("&v: error communicating to the spark cloud", resp.Status))
		}
		return
	}

	return
}
