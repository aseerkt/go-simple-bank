package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {

	testCases := []struct {
		name         string
		secretKey    string
		checkResults func(Maker, error)
	}{
		{
			name:      "Ok",
			secretKey: "secretsthatgoespublicareeviltruth",
			checkResults: func(m Maker, e error) {
				require.NoError(t, e)
				require.NotEmpty(t, m)
			},
		},
		{
			name:      "InvalidSecretKey",
			secretKey: "tooshort",
			checkResults: func(m Maker, e error) {
				require.Error(t, e)
				require.Empty(t, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checkResults(NewJWTMaker(tc.secretKey))
		})
	}

}

func TestCreateToken(t *testing.T) {
	jm, err := NewJWTMaker("thelongestsecretkeythatyouwillneverimagine")

	require.NoError(t, err)

	token, err := jm.CreateToken("armageddon", 12*time.Hour)

	require.NoError(t, err)
	require.NotZero(t, token)
}

func TestVerifyToken(t *testing.T) {

	secretKey := "thesecretjustforverifyingsomeshit"

	jm, err := NewJWTMaker(secretKey)

	require.NoError(t, err)

	testCases := []struct {
		name         string
		createToken  func() string
		checkResults func(*Payload, error)
	}{
		{
			name: "Ok",
			createToken: func() string {
				token, err := jm.CreateToken("alfred", 12*time.Hour)

				require.NoError(t, err)
				return token
			},
			checkResults: func(p *Payload, err error) {
				require.NoError(t, err)
				require.Equal(t, p.Username, "alfred")
				require.NotZero(t, p.ID)
				require.NotZero(t, p.IssuedAt)
				require.NotZero(t, p.ExpiresAt)
			},
		},
		{
			name: "InvalidSignature",
			createToken: func() string {
				jwtToken := jwt.New(jwt.SigningMethodES256)
				key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t, err)
				token, err := jwtToken.SignedString(key)
				require.NoError(t, err)
				return token

			},
			checkResults: func(p *Payload, err error) {
				require.Empty(t, p)
				require.Error(t, err)
				require.ErrorIs(t, err, jwt.ErrTokenUnverifiable)
			},
		},
		{
			name: "InvalidClaims",
			createToken: func() string {
				jwtToken := jwt.New(jwt.SigningMethodHS256)
				token, err := jwtToken.SignedString([]byte(secretKey))
				require.NoError(t, err)
				return token
			},
			checkResults: func(p *Payload, err error) {
				require.Empty(t, p)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := tc.createToken()
			p, err := jm.VerifyToken(token)
			tc.checkResults(p, err)
		})
	}
}
