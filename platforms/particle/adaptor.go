package particle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donovanhide/eventsource"
	"gobot.io/x/gobot"
)

// Adaptor is the Gobot Adaptor for Particle
type Adaptor struct {
	name        string
	DeviceID    string
	AccessToken string
	APIServer   string
	servoPins   map[string]bool
	gobot.Eventer
}

// Event is an event emitted by the Particle cloud
type Event struct {
	Name  string
	Data  string
	Error error
}

var eventSource = func(url string) (chan eventsource.Event, chan error, error) {
	stream, err := eventsource.Subscribe(url, "")
	if err != nil {
		return nil, nil, err
	}
	return stream.Events, stream.Errors, nil
}

// NewAdaptor creates new Photon adaptor with deviceId and accessToken
// using api.particle.io server as default
func NewAdaptor(deviceID string, accessToken string) *Adaptor {
	return &Adaptor{
		name:        gobot.DefaultName("Particle"),
		DeviceID:    deviceID,
		AccessToken: accessToken,
		servoPins:   make(map[string]bool),
		APIServer:   "https://api.particle.io",
		Eventer:     gobot.NewEventer(),
	}
}

// Name returns the Adaptor name
func (s *Adaptor) Name() string { return s.name }

// SetName sets the Adaptor name
func (s *Adaptor) SetName(n string) { s.name = n }

// Connect returns true if connection to Particle Photon or Electron is successful
func (s *Adaptor) Connect() (err error) {
	return
}

// Finalize returns true if connection to Particle Photon or Electron is finalized successfully
func (s *Adaptor) Finalize() (err error) {
	return
}

// AnalogRead reads analog ping value using Particle cloud api
func (s *Adaptor) AnalogRead(pin string) (val int, err error) {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}

	url := fmt.Sprintf("%v/analogread", s.deviceURL())

	resp, err := s.request("POST", url, params)
	if err == nil {
		val = int(resp["return_value"].(float64))
		return
	}

	return 0, err
}

// PwmWrite writes in pin using analog write api
func (s *Adaptor) PwmWrite(pin string, level byte) (err error) {
	return s.AnalogWrite(pin, level)
}

// AnalogWrite writes analog pin with specified level using Particle cloud api
func (s *Adaptor) AnalogWrite(pin string, level byte) (err error) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, level)},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/analogwrite", s.deviceURL())
	_, err = s.request("POST", url, params)
	return
}

// DigitalWrite writes to a digital pin using Particle cloud api
func (s *Adaptor) DigitalWrite(pin string, level byte) (err error) {
	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, s.pinLevel(level))},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalwrite", s.deviceURL())
	_, err = s.request("POST", url, params)
	return err
}

// DigitalRead reads from digital pin using Particle cloud api
func (s *Adaptor) DigitalRead(pin string) (val int, err error) {
	params := url.Values{
		"params":       {pin},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/digitalread", s.deviceURL())
	resp, err := s.request("POST", url, params)
	if err == nil {
		val = int(resp["return_value"].(float64))
		return
	}
	return -1, err
}

// ServoWrite writes the 0-180 degree angle to the specified pin.
// To use it requires installing the "tinker-servo" sketch on your
// Particle device. not just the default "tinker".
func (s *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	if _, present := s.servoPins[pin]; !present {
		err = s.servoPinOpen(pin)
		if err != nil {
			return
		}
	}

	params := url.Values{
		"params":       {fmt.Sprintf("%v,%v", pin, angle)},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/servoSet", s.deviceURL())
	_, err = s.request("POST", url, params)
	return err
}

// EventStream returns a gobot.Event based on the following params:
//
// * source - "all"/"devices"/"device" (More info at: http://docs.particle.io/api/#reading-data-from-a-core-events)
// * name  - Event name to subscribe for, leave blank to subscribe to all events.
//
// A new event is emitted as a particle.Event struct
func (s *Adaptor) EventStream(source string, name string) (event *gobot.Event, err error) {
	var url string

	switch source {
	case "all":
		url = fmt.Sprintf("%s/v1/events/%s?access_token=%s", s.APIServer, name, s.AccessToken)
	case "devices":
		url = fmt.Sprintf("%s/v1/devices/events/%s?access_token=%s", s.APIServer, name, s.AccessToken)
	case "device":
		url = fmt.Sprintf("%s/events/%s?access_token=%s", s.deviceURL(), name, s.AccessToken)
	default:
		err = errors.New("source param should be: all, devices or device")
		return
	}

	events, _, err := eventSource(url)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case ev := <-events:
				if ev.Event() != "" && ev.Data() != "" {
					s.Publish(ev.Event(), ev.Data())
				}
			}
		}
	}()
	return
}

// Variable returns a core variable value as a string
func (s *Adaptor) Variable(name string) (result string, err error) {
	url := fmt.Sprintf("%v/%s?access_token=%s", s.deviceURL(), name, s.AccessToken)
	resp, err := s.request("GET", url, nil)

	if err != nil {
		return
	}

	val := resp["result"]
	switch val.(type) {
	case bool:
		result = strconv.FormatBool(val.(bool))
	case float64:
		result = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case string:
		result = val.(string)
	}

	return
}

// Function executes a core function and
// returns value from request.
// Takes a String as the only argument and returns an Int.
// If function is not defined in core, it will time out
func (s *Adaptor) Function(name string, args string) (val int, err error) {
	params := url.Values{
		"args":         {args},
		"access_token": {s.AccessToken},
	}

	url := fmt.Sprintf("%s/%s", s.deviceURL(), name)
	resp, err := s.request("POST", url, params)

	if err != nil {
		return -1, err
	}

	val = int(resp["return_value"].(float64))
	return
}

// setAPIServer sets Particle cloud api server, this can be used to change from default api.spark.io
func (s *Adaptor) setAPIServer(server string) {
	s.APIServer = server
}

// deviceURL constructs device url to make requests from Particle cloud api
func (s *Adaptor) deviceURL() string {
	if len(s.APIServer) <= 0 {
		s.setAPIServer("https://api.particle.io")
	}
	return fmt.Sprintf("%v/v1/devices/%v", s.APIServer, s.DeviceID)
}

// pinLevel converts byte level to string expected in api
func (s *Adaptor) pinLevel(level byte) string {
	if level == 1 {
		return "HIGH"
	}
	return "LOW"
}

// request makes request to Particle cloud server, return err != nil if there is
// any issue with the request.
func (s *Adaptor) request(method string, url string, params url.Values) (m map[string]interface{}, err error) {
	var resp *http.Response

	if method == "POST" {
		resp, err = http.PostForm(url, params)
	} else if method == "GET" {
		resp, err = http.Get(url)
	}

	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	json.Unmarshal(buf, &m)

	if resp.Status != "200 OK" {
		err = fmt.Errorf("%v: error communicating to the Particle cloud", resp.Status)
	} else if _, ok := m["error"]; ok {
		err = errors.New(m["error"].(string))
	}

	return
}

func (s *Adaptor) servoPinOpen(pin string) error {
	params := url.Values{
		"params":       {fmt.Sprintf("%v", pin)},
		"access_token": {s.AccessToken},
	}
	url := fmt.Sprintf("%v/servoOpen", s.deviceURL())
	_, err := s.request("POST", url, params)
	if err != nil {
		return err
	}
	s.servoPins[pin] = true
	return nil
}
