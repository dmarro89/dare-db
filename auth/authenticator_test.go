package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	authenticator := NewJWTAutenticator()
	username := "testuser"

	token, err := authenticator.GenerateToken(username)
	assert.NoError(t, err, "Expected no error when generating token")
	assert.NotEmpty(t, token, "Expected a token, got an empty string")

	// Optional: Parse the token to verify it's correctly formatted
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	assert.NoError(t, err, "Error parsing token")
	assert.True(t, parsedToken.Valid, "Generated token is not valid")
	assert.Equal(t, username, claims.Username, "Expected username %v, got %v", username, claims.Username)
}

func TestJWTAuthenticator_VerifyToken(t *testing.T) {
	userStore := NewUserStore()
	authenticator := &JWTAutenticator{
		usersStore: userStore,
	}

	username := "testuser"
	password := "testpassword"
	userStore.AddUser(username, password)
	tokenString, err := authenticator.GenerateToken(username)
	assert.NoError(t, err, "Expected no error when generating token")
	assert.NotEmpty(t, tokenString, "Expected a token, got an empty string")

	userStore.SaveToken(username, tokenString)

	// Test case: valid token
	returnedUsername, err := authenticator.VerifyToken(tokenString)
	assert.NoError(t, err, "Expected no error for a valid token")
	assert.Equal(t, username, returnedUsername, "Expected username to match")

	// Test case: invalid token
	invalidTokenString := tokenString + "invalid"
	_, err = authenticator.VerifyToken(invalidTokenString)
	assert.Error(t, err, "Expected error for an invalid token")
	assert.Equal(t, "failed to parse token: token signature is invalid: signature is invalid", err.Error(), "Expected 'token signature is invalid: signature is invalid' error")

	// Test case: expired token
	expiredTokenString, _ := generateExpiredToken(username)
	userStore.SaveToken(username, expiredTokenString)
	_, err = authenticator.VerifyToken(expiredTokenString)
	assert.Error(t, err, "Expected error for an expired token")

	// Test case: token not present in user store
	userStore.DeleteUser(username)
	_, err = authenticator.VerifyToken(tokenString)
	assert.Error(t, err, "Expected error when the token is not present in the user store")
	assert.Equal(t, "invalid token", err.Error(), "Expected 'invalid token' error")
}

// Helper function to generate an expired token
func generateExpiredToken(username string) (string, error) {
	expirationTime := time.Now().Add(-5 * time.Minute) // Set expiration time in the past
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
