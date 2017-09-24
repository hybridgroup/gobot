package api

import (
	"net/http"
	"regexp"
	"strings"
)

// CORS represents CORS configuration
type CORS struct {
	AllowOrigins        []string
	AllowHeaders        []string
	AllowMethods        []string
	ContentType         string
	allowOriginPatterns []string
}

// AllowRequestsFrom returns handler to verify that requests come from allowedOrigins
func AllowRequestsFrom(allowedOrigins ...string) http.HandlerFunc {
	c := &CORS{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		ContentType:  "application/json; charset=utf-8",
	}

	c.generatePatterns()

	return func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if c.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", c.AllowedHeaders())
			w.Header().Set("Access-Control-Allow-Methods", c.AllowedMethods())
			w.Header().Set("Content-Type", c.ContentType)
		}
	}
}

// isOriginAllowed returns true if origin matches an allowed origin pattern.
func (c *CORS) isOriginAllowed(origin string) (allowed bool) {
	for _, allowedOriginPattern := range c.allowOriginPatterns {
		allowed, _ = regexp.MatchString(allowedOriginPattern, origin)
		if allowed {
			return
		}
	}
	return
}

// generatePatterns generates regex expression for AllowOrigins
func (c *CORS) generatePatterns() {
	if c.AllowOrigins != nil {
		for _, origin := range c.AllowOrigins {
			pattern := regexp.QuoteMeta(origin)
			pattern = strings.Replace(pattern, "\\*", ".*", -1)
			pattern = strings.Replace(pattern, "\\?", ".", -1)
			c.allowOriginPatterns = append(c.allowOriginPatterns, "^"+pattern+"$")
		}
	}
}

// AllowedHeaders returns allowed headers in a string
func (c *CORS) AllowedHeaders() string {
	return strings.Join(c.AllowHeaders, ",")
}

// AllowedMethods returns allowed http methods in a string
func (c *CORS) AllowedMethods() string {
	return strings.Join(c.AllowMethods, ",")
}
