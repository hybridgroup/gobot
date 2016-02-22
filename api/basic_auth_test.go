package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestBasicAuth(t *testing.T) {
	a := initTestAPI()

	a.AddHandler(BasicAuth("admin", "password"))

	request, _ := http.NewRequest("GET", "/api/", nil)
	request.SetBasicAuth("admin", "password")
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)
	gobottest.Assert(t, response.Code, 200)

	request, _ = http.NewRequest("GET", "/api/", nil)
	request.SetBasicAuth("admin", "wrongPassword")
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	gobottest.Assert(t, response.Code, 401)
}
