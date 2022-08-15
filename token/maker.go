package token

import "time"

type Maker interface {
	CreateToken(userID string, duration time.Duration) (string, error)

	ValidateToken(token string) (*Payload, error)
}
