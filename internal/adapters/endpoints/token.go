package endpoints

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

// UserClaim jwt encoded object
type UserClaim struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

func GetTokenFromAuthorization(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

// DecodeUserFromToken decodes UserClaim from the given jwt
func DecodeUserFromToken(token string) (*UserClaim, error) {
	claim, err := jwt.ParseWithClaims(token, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return claim.Claims.(*UserClaim), nil
}

// GetUserClaimFromRequest get decoded UserClaim from the r Request
func GetUserClaimFromRequest(r *http.Request) (*UserClaim, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("unauthorized")
	}
	authToken := GetTokenFromAuthorization(authHeader)
	userClaim, err := DecodeUserFromToken(authToken)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return userClaim, nil
}
