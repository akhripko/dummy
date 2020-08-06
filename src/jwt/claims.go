package jwt

import (
	"github.com/akhripko/dummy/src/models"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserID      string             `json:"userId"`
	Account     string             `json:"account"`
	Permissions models.Permissions `json:"permissions"`
	jwtgo.StandardClaims
}
