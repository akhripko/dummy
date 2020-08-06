package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
)

func Test_BuildTokenFlow(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := TokenClaims{
		UserID:  "user id",
		Account: "acc name",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(rsaPrivateKey)
	assert.NoError(t, err)

	publicKeyBytes, err := ExportRsaPublicKeyAsPemStr(&rsaPrivateKey.PublicKey)
	assert.NoError(t, err)

	var claims2 TokenClaims
	tkn, err := jwt.ParseWithClaims(tokenString, &claims2, func(token *jwt.Token) (interface{}, error) {
		return ParseRsaPublicKeyFromPemStr(publicKeyBytes)
	})
	assert.NoError(t, err)
	assert.True(t, tkn.Valid)
	assert.Equal(t, claims, claims2)
}

func Test_BuildParseFlow(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := TokenClaims{
		UserID:  "user id",
		Account: "acc id",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	tokenString, err := BuildSignedToken(rsaPrivateKey, "keyID", claims)
	assert.NoError(t, err)

	parsedClaims, err := ParseSignedToken(tokenString, &rsaPrivateKey.PublicKey)
	assert.NoError(t, err)
	assert.NoError(t, parsedClaims.Valid())

	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	assert.Equal(t, claims.Account, parsedClaims.Account)
}
