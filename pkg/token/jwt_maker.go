package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeyLength = 12

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid secret key size: minimum %d characters required", minSecretKeyLength)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (jm *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwt.SignedString([]byte(jm.secretKey))

}

func (jm *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var keyFunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {

		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(jm.secretKey), nil
	}

	jwt, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if payload, ok := jwt.Claims.(*Payload); !ok {
		return nil, errors.New("invalid token")
	} else {
		return payload, nil
	}

}
