package auth

import "time"

const (
	JwT_ISSUER       = "altar"
	JWT_ACCESS_TOKEN_TTL = 15 * time.Minute
)
