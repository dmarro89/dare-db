package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dmarro89/dare-db/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	modelFile, err := os.CreateTemp("", "rbac_model.conf")
	if err != nil {
		t.Fatalf("Error creating rbac model file: %v", err)
	}
	defer os.Remove(modelFile.Name())

	policyFile, err := os.CreateTemp("", "rbac_policy.csv")
	if err != nil {
		t.Fatalf("Error creating rbac policy file: %v", err)
	}
	defer os.Remove(policyFile.Name())

	if _, err := modelFile.Write([]byte(RBAC_MODEL_CONTENT)); err != nil {
		t.Fatalf("Errorwriting rbac model file: %v", err)
	}
	if err := modelFile.Close(); err != nil {
		t.Fatalf("Error closing rbac model file: %v", err)
	}

	if _, err := policyFile.Write([]byte(RBAC_POLICY)); err != nil {
		t.Fatalf("Error creating policy file: %v", err)
	}
	if err := policyFile.Close(); err != nil {
		t.Fatalf("Error closing policy file: %v", err)
	}

	casbinAuth := NewCasbinAuth(modelFile.Name(), policyFile.Name(), Users{
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
	require.NoError(t, err)
	assert.NotNil(t, token)
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
