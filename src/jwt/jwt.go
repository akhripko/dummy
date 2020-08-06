package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func BuildSignedToken(privateKey *rsa.PrivateKey, keyID string, claims TokenClaims) (string, error) {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	// Create the JWT string
	return token.SignedString(privateKey)
}

func ParseSignedToken(tokenString string, publicKey *rsa.PublicKey) (*TokenClaims, error) {
	var claims TokenClaims
	_, err := jwtgo.ParseWithClaims(tokenString, &claims, func(token *jwtgo.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id token string")
	}
	return &claims, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}
