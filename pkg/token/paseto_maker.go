package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoMaker() Maker {
	return &PasetoMaker{paseto.NewV4SymmetricKey(), []byte("some_implicit")}
}

func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	token.SetString("id", uuid.New().String())
	token.SetString("username", username)
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	return token.V4Encrypt(pm.symmetricKey, pm.implicit), nil
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(pm.symmetricKey, token, pm.implicit)

	if err != nil {
		return nil, err
	}
	payload, err := getPayloadFromToken(parsedToken)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func getPayloadFromToken(token *paseto.Token) (*Payload, error) {
	id, err := token.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	username, err := token.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiresAt, err := token.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}, nil
}
