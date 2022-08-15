package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secret key is too short, must be at least %d bytes", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userID string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) ValidateToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*Payload); ok && jwtToken.Valid {
		return claims, nil
	} else {
		return nil, ErrorInvalidToken
	}
}
