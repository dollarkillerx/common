package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// AccessClaims ...
type AccessClaims struct {
	Info Info `json:"info"`
	jwt.StandardClaims
}

type Info struct {
	Info map[string]string `json:"info"`
}

// CreateAccessToken ...
func CreateAccessToken(Info Info, secretKey string) (string, error) {
	claims := AccessClaims{
		Info,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}
	return CreateTokenString(secretKey, claims)
}

// GetInfoFromAccessToken ...
func GetInfoFromAccessToken(tokenString string, secretKey string) (Info, error) {
	token, err := ParseTokenString(tokenString, secretKey, &AccessClaims{})
	if err != nil {
		return Info{}, err
	}
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims.Info, nil
	}
	return Info{}, errors.New("token error")
}

// GetClaimsFromAccessToken ...
func GetClaimsFromAccessToken(tokenString string, secretKey string) (AccessClaims, error) {
	token, err := ParseTokenString(tokenString, secretKey, &AccessClaims{})
	if err != nil {
		return AccessClaims{}, err
	}
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return *claims, nil
	}
	return AccessClaims{}, errors.New("token error")
}

// CreateTokenString ...
func CreateTokenString(secretKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return tokenString, nil
}

// ParseTokenString ...
func ParseTokenString(tokenString, secretKey string, claims jwt.Claims) (*jwt.Token, error) {
	v, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return v, nil
}
