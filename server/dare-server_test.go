package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_SetAndGet(t *testing.T) {
	db := database.NewDatabase()
	srv := NewDareServer(db, auth.NewUserStore())

	setWrongResponse := httptest.NewRecorder()
	setWrongRequest, _ := http.NewRequest("GET", "/set", bytes.NewBuffer([]byte{}))

	srv.HandlerSet(setWrongResponse, setWrongRequest)
	assert.Equal(t, http.StatusMethodNotAllowed, setWrongResponse.Code, "Method not allowed")

	setWrongFormatRequest, _ := http.NewRequest("POST", "/set", bytes.NewBuffer([]byte("plainText")))
	setWrongFormatResponse := httptest.NewRecorder()
	srv.HandlerSet(setWrongFormatResponse, setWrongFormatRequest)
	assert.Equal(t, http.StatusBadRequest, setWrongFormatResponse.Code, "Invalid JSON format, the body must be in the form of {\"key\": \"value\"}")

	setData := map[string]interface{}{"testKey": "testValue"}
	setDataJSON, _ := json.Marshal(setData)
	setRequest, _ := http.NewRequest("POST", "/set", bytes.NewBuffer(setDataJSON))

	setResponse := httptest.NewRecorder()
	srv.HandlerSet(setResponse, setRequest)
	assert.Equal(t, http.StatusCreated, setResponse.Code, "Expected status %d, got %d", http.StatusCreated, setResponse.Code)

	getEmptyRequest, _ := http.NewRequest("GET", "/get", nil)
	getEmptyResponse := httptest.NewRecorder()
	srv.HandlerGetById(getEmptyResponse, getEmptyRequest)
	assert.Equal(t, http.StatusBadRequest, getEmptyResponse.Code, `url param query "key" cannot be empty`)

	getMissingKeyRequest, _ := http.NewRequest("GET", "/get/missingKey", nil)
	getMissingKeyRequest.SetPathValue("key", "missingKey")
	getMissingKeyResponse := httptest.NewRecorder()
	srv.HandlerGetById(getMissingKeyResponse, getMissingKeyRequest)
	assert.Equal(t, http.StatusNotFound, getMissingKeyResponse.Code, `Key "%s" not found`, "missingKey")

	getRequest, _ := http.NewRequest("GET", "/get/testKey", nil)
	getRequest.SetPathValue("key", "testKey")
	getResponse := httptest.NewRecorder()

	srv.HandlerGetById(getResponse, getRequest)

	assert.Equal(t, http.StatusOK, getResponse.Code, "Expected status %d, got %d", http.StatusOK, getResponse.Code)

	var getResult map[string]interface{}
	err := json.Unmarshal(getResponse.Body.Bytes(), &getResult)
	assert.Nil(t, err, "Error decoding JSON response")

	for key, value := range getResult {
		assert.Equal(t, key, "testKey", "Unexpected response body: %v", getResult)
		assert.Equal(t, value, "testValue", "Unexpected response body: %v", getResult)
	}
}

func TestServer_SetAndDelete(t *testing.T) {
	db := database.NewDatabase()
	srv := NewDareServer(db, auth.NewUserStore())

	setData := map[string]interface{}{"testKey": "testValue"}
	setDataJSON, _ := json.Marshal(setData)
	setRequest, _ := http.NewRequest("POST", "/set", bytes.NewBuffer(setDataJSON))
	setResponse := httptest.NewRecorder()

	srv.HandlerSet(setResponse, setRequest)

	assert.Equal(t, http.StatusCreated, setResponse.Code, "Expected status %d, got %d", http.StatusCreated, setResponse.Code)

	deleteWrongResponse := httptest.NewRecorder()
	deleteWrongRequest, _ := http.NewRequest("GET", "/delete", nil)
	srv.HandlerDelete(deleteWrongResponse, deleteWrongRequest)
	assert.Equal(t, http.StatusMethodNotAllowed, deleteWrongResponse.Code, "Method not allowed")

	deleteEmptyRequest, _ := http.NewRequest("DELETE", "/delete", nil)
	deleteEmptyResponse := httptest.NewRecorder()
	srv.HandlerDelete(deleteEmptyResponse, deleteEmptyRequest)
	assert.Equal(t, http.StatusBadRequest, deleteEmptyResponse.Code, `url param query "key" cannot be empty`)

	deleteRequest, _ := http.NewRequest("DELETE", "/delete/testKey", nil)
	deleteRequest.SetPathValue("key", "testKey")
	deleteResponse := httptest.NewRecorder()

	srv.HandlerDelete(deleteResponse, deleteRequest)

	if deleteResponse.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, deleteResponse.Code)
	}

	getRequest, _ := http.NewRequest("GET", "/get/testKey", nil)
	getRequest.SetPathValue("key", "testKey")
	getResponse := httptest.NewRecorder()

	srv.HandlerGetById(getResponse, getRequest)

	if getResponse.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, getResponse.Code)
	}
}

func TestCreateMux(t *testing.T) {
	// Create a new instance of DareServer
	db := database.NewDatabase()
	srv := NewDareServer(db, auth.NewUserStore())

	// Create a new ServeMux using the CreateMux method
	mux := srv.CreateMux(auth.NewCasbinAuth("../auth/test/rbac_model.conf", "../auth/test/rbac_policy.csv", auth.Users{
		"user1": {Roles: []string{"role1"}}, "user2": {Roles: []string{"role2"}},
	}), auth.NewJWTAutenticator())
	assert.Equal(t, mux != nil, true)
}

func TestMiddleware_ProtectedEndpoints(t *testing.T) {
	db := database.NewDatabase()

	usersStore := auth.NewUserStore()
	usersStore.AddUser("user1", "password")
	usersStore.AddUser("user2", "password")

	srv := NewDareServer(db, usersStore)

	authenticator := auth.NewJWTAutenticatorWithUsers(usersStore)
	mux := srv.CreateMux(auth.NewCasbinAuth("../auth/test/rbac_model.conf", "../auth/test/rbac_policy.csv", auth.Users{
		"user1": {Roles: []string{"role1"}}, "user2": {Roles: []string{"role2"}},
	}), authenticator)

	req, err := http.NewRequest("POST", "/login", nil)
	require.NoError(t, err)
	req.SetBasicAuth("user2", "password")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var tokenResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&tokenResponse)

	// Test POST request for a protected resource
	postData := map[string]interface{}{"newKey": "newValue"}
	body, err := json.Marshal(postData)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "/set", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", tokenResponse["token"])

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check if the response status code is Created (201)
	require.Equal(t, http.StatusCreated, rr.Code)

	//Login user1
	req, err = http.NewRequest("POST", "/login", nil)
	require.NoError(t, err)
	req.SetBasicAuth("user1", "password")

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	json.NewDecoder(rr.Body).Decode(&tokenResponse)

	// Test GET request for a protected resource
	req, err = http.NewRequest("GET", "/get/newKey", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", tokenResponse["token"])

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check if the response status code is OK (200)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Body.String(), string(body))

	// Test DELETE request for a protected resource
	req, err = http.NewRequest("DELETE", "/delete/existingKey", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", tokenResponse["token"])

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check if the response status code is OK (200)
	require.Equal(t, http.StatusForbidden, rr.Code)

	// Test accessing without proper credentials
	req, err = http.NewRequest("GET", "/get/existingKey", nil)
	require.NoError(t, err)
	// Not setting basic auth here to simulate missing credentials

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check if the response status code is Unauthorized (401)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDareServer_HandlerLogin(t *testing.T) {
	usersStore := auth.NewUserStore()
	server := &DareServer{
		usersStore: usersStore,
	}

	// Adding a test user
	usersStore.AddUser("testuser", "testpassword")

	// Test case: valid login request
	req, err := http.NewRequest(http.MethodPost, "/login", nil)
	assert.NoError(t, err)

	req.SetBasicAuth("testuser", "testpassword")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.HandlerLogin)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]string
	err = json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "token")
	assert.NotEmpty(t, responseBody["token"])

	// Test case: invalid method
	req, err = http.NewRequest(http.MethodGet, "/login", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// Test case: missing credentials
	req, err = http.NewRequest(http.MethodPost, "/login", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Test case: invalid credentials
	req, err = http.NewRequest(http.MethodPost, "/login", nil)
	assert.NoError(t, err)

	req.SetBasicAuth("testuser", "wrongpassword")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
