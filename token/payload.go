package token

import (
	"errors"
	"time"

	"github.com/go-oauth2/oauth2/utils/uuid"
)

var ErrorExpiredToken = errors.New("expired token")
var ErrorInvalidToken = errors.New("invalid token")

type Payload struct {
	PayloadID string    `json:"payload_id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		println("Error at creating payload", err)
		return nil, err
	}

	return &Payload{
		PayloadID: id.String(),
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
