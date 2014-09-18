package api

import (
	"regexp"
	"strings"
)

type CORS struct {
	AllowOrigins        []string
	AllowHeaders        []string
	AllowMethods        []string
	ContentType         string
	allowOriginPatterns []string
}

func NewCORS(allowedOrigins []string) *CORS {
	cors := &CORS{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		ContentType:  "application/json; charset=utf-8",
	}

	cors.generatePatterns()

	return cors
}

func (c *CORS) isOriginAllowed(origin string) (allowed bool) {
	for _, allowedOriginPattern := range c.allowOriginPatterns {
		allowed, _ = regexp.MatchString(allowedOriginPattern, origin)
		if allowed {
			return
		}
	}
	return
}

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

func (c *CORS) AllowedHeaders() string {
	return strings.Join(c.AllowHeaders, ",")
}

func (c *CORS) AllowedMethods() string {
	return strings.Join(c.AllowMethods, ",")
}
