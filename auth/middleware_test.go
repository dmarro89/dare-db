package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmarro89/dare-db/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	casbinAuth := NewCasbinAuth(modelPath, policyPath, Users{
		"user1": {Roles: []string{"role1"}},
		"user2": {Roles: []string{"role2"}},
	})

	userStore := NewUserStore()
	authenticator := &JWTAutenticator{
		usersStore: userStore,
	}

	middleware := &DareMiddleware{
		authorizer:    casbinAuth,
		authenticator: authenticator,
		logger:        logger.NewDareLogger(),
	}

	handler := middleware.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test case where user1 is authorized for GET
	req, err := http.NewRequest("GET", "/some-path", nil)
	require.NoError(t, err)

	token, err := middleware.authenticator.GenerateToken("user1")
	assert.Nil(t, err)
	userStore.SaveToken("user1", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	// Test case where user1 is not authorized for POST
	req, err = http.NewRequest("POST", "/some-path", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", token)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)

	// Test case where user2 is authorized for POST
	req, err = http.NewRequest("POST", "/some-path", nil)
	require.NoError(t, err)
	token, err = middleware.authenticator.GenerateToken("user2")
	assert.Nil(t, err)
	userStore.SaveToken("user2", token)
	req.Header.Set("Authorization", token)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	// Test case where user2 is not authorized for GET
	req, err = http.NewRequest("GET", "/some-path", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", token)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)

	// Test case where credentials are missing
	req, err = http.NewRequest("GET", "/some-path", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
