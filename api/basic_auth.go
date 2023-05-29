package api

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
)

// BasicAuth returns basic auth handler.
func BasicAuth(username, password string) http.HandlerFunc {
	// Inspired by https://github.com/codegangsta/martini-contrib/blob/master/auth/
	return func(res http.ResponseWriter, req *http.Request) {
		if !secureCompare(req.Header.Get("Authorization"),
			"Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
		) {
			res.Header().Set("WWW-Authenticate",
				"Basic realm=\"Authorization Required\"",
			)
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}

func secureCompare(given string, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	}
	// Securely compare actual to itself to keep constant time,
	// but always return false
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
}
