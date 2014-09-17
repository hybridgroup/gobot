package api

type CORS struct {
	AllowOrigins []string
	AllowHeaders []string // Not yet implemented
	AllowMethods []string // ditto
}

func NewCORS(allowedOrigins []string) *CORS {
	return &CORS{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}
}

func (c *CORS) isOriginAllowed(currentOrigin string) bool {
	for _, allowedOrigin := range c.AllowOrigins {
		if "*" == allowedOrigin || currentOrigin == allowedOrigin {
			return true
		}
	}
	return false
}
