package dronesmith

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hybridgroup/gobot"
)

type Adaptor struct {
	name      string
	DroneID   string
	UserEmail string
	UserKey   string
	APIServer string
	gobot.Eventer
}

// NewAdaptor creates new Dronesmith adaptor with droneId and accessToken
// using api.dronesmith.io server as default
func NewAdaptor(droneID string, userEmail string, userKey string) *Adaptor {
	return &Adaptor{
		name:      "Dronesmith",
		DroneID:   droneID,
		UserEmail: userEmail,
		UserKey:   userKey,
		APIServer: "http://api.dronesmith.io",
		Eventer:   gobot.NewEventer(),
	}
}

func (s *Adaptor) Name() string     { return s.name }
func (s *Adaptor) SetName(n string) { s.name = n }

// Connect returns nil if connection to Dronesmith server is successful,
// otherwise returns the http error
func (s *Adaptor) Connect() (err error) {
	_, err = s.Request("POST", "/start", nil)

	return
}

// Finalize returns nil if connection to Dronesmith server is finalized successfully
// otherwise returns the error
func (s *Adaptor) Finalize() (err error) {
	_, err = s.Request("POST", "/stop", nil)

	return
}

// SetAPIServer sets Dronesmith api server, this can be used to change from default api.dronesmith.io
func (s *Adaptor) setAPIServer(server string) {
	s.APIServer = server
}

// droneURL constructs drone url to make requests from Dronesmith api
func (s *Adaptor) droneURL() string {
	if len(s.APIServer) <= 0 {
		s.setAPIServer("http://api.dronesmith.io")
	}
	return fmt.Sprintf("%v/api/drone/%v", s.APIServer, s.DroneID)
}

func (s *Adaptor) Request(method string, url string, params url.Values) (m map[string]interface{}, err error) {
	fullURL := fmt.Sprintf("%v%v", s.droneURL(), url)
	m, err = s.request(method, fullURL, params)
	return
}

// request makes request to Dronesmith server, return err != nil if there is
// any issue with the request.
func (s *Adaptor) request(method string, url string, params url.Values) (m map[string]interface{}, err error) {
	client := &http.Client{}
	var req *http.Request
	var resp *http.Response

	if method == "POST" {
		req, _ = http.NewRequest(method, url, strings.NewReader(params.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if method == "GET" {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Set("user-email", s.UserEmail)
	req.Header.Set("user-key", s.UserKey)

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	json.Unmarshal(buf, &m)

	if resp.Status != "200 OK" {
		err = fmt.Errorf("%v: error communicating to the Dronesmith server", resp.Status)
	}

	return
}
