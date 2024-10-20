package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// constants
const JWT_TIME_TO_LIVE_MINUTES int = 60

type Authenticator interface {
	GenerateToken(string) (string, error)
	VerifyToken(token string) (string, error)
}

type JWTAutenticator struct {
	usersStore *UserStore
	jwtKey     []byte
}

func NewJWTAutenticator() *JWTAutenticator {
	return &JWTAutenticator{
		jwtKey:     getJWTKey(),
		usersStore: NewUserStore(),
	}
}

func NewJWTAutenticatorWithUsers(usersStore *UserStore) *JWTAutenticator {
	return &JWTAutenticator{
		jwtKey:     getJWTKey(),
		usersStore: usersStore,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (jwtAuthenticator *JWTAutenticator) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(JWT_TIME_TO_LIVE_MINUTES) * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAuthenticator.jwtKey)
}

func (jwtAuthenticator *JWTAutenticator) VerifyToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtAuthenticator.jwtKey, nil
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

var (
	jwtKey []byte
	once   sync.Once
)

func getJWTKey() []byte {
	once.Do(func() {
		key := os.Getenv("JWT_SECRET_KEY")

		if key == "" {
			randomKey := make([]byte, 32)
			if _, err := rand.Read(randomKey); err != nil {
				panic("Failed to generate random key: " + err.Error())
			}
			key = hex.EncodeToString(randomKey)
		}

		jwtKey = []byte(key)
	})

	return jwtKey
}
