//nolint:usestdlibvars,noctx // ok here
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	a := initTestAPI()

	a.AddHandler(BasicAuth("admin", "password"))

	request, _ := http.NewRequest("GET", "/api/", nil)
	request.SetBasicAuth("admin", "password")
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	request, _ = http.NewRequest("GET", "/api/", nil)
	request.SetBasicAuth("admin", "wrongPassword")
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 401, response.Code)
}
