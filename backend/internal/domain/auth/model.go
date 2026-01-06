package auth

import "time"

// TokenPair represents access and refresh tokens issued to a client.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	AccessExp    time.Time
	RefreshExp   time.Time
}
