package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateToken(string) (string, error)
	VerifyToken(token string) (string, error)
}

type JWTAutenticator struct {
	usersStore *UserStore
}

func NewJWTAutenticator() *JWTAutenticator {
	return &JWTAutenticator{
		usersStore: NewUserStore(),
	}
}

func NewJWTAutenticatorWithUsers(usersStore *UserStore) *JWTAutenticator {
	return &JWTAutenticator{
		usersStore: usersStore,
	}
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (jwtAuthenticator *JWTAutenticator) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func (jwtAuthenticator *JWTAutenticator) VerifyToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userStore := jwtAuthenticator.usersStore

	if userStore.tokens[claims.Username] != tokenString {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}
