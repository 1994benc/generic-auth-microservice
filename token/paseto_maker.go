package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// CreateToken implements Maker
func (maker *PasetoMaker) CreateToken(userID string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// ValidateToken implements Maker
func (maker *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewPasetoMaker(symmetricKey []byte) (Maker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symmetricKey must be %d bytes long", chacha20poly1305.KeySize)
	}

	paseto := paseto.NewV2()

	return &PasetoMaker{
		paseto:       paseto,
		symmetricKey: []byte(symmetricKey),
	}, nil
}
