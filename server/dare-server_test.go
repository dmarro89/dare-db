package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

	setData := map[string]string{"testKey": "testValue"}
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

	var getResult map[string]string
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

	setData := map[string]string{"testKey": "testValue"}
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

const RBAC_MODEL_CONTENT = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (p.obj == "*" || keyMatch(r.obj, p.obj)) && regexMatch(r.act, p.act)
`

const RBAC_POLICY = `p, role1, *, GET
p, role2, *, POST

g, user1, role1
g, user2, role2
`

func TestCreateMux(t *testing.T) {
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

	// Create a new instance of DareServer
	db := database.NewDatabase()
	srv := NewDareServer(db, auth.NewUserStore())

	// Create a new ServeMux using the CreateMux method
	mux := srv.CreateMux(auth.NewCasbinAuth(modelFile.Name(), policyFile.Name(), auth.Users{
		"user1": {Roles: []string{"role1"}}, "user2": {Roles: []string{"role2"}},
	}), auth.NewJWTAutenticator())
	assert.Equal(t, mux != nil, true)
}

func TestMiddleware_ProtectedEndpoints(t *testing.T) {
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

	db := database.NewDatabase()

	usersStore := auth.NewUserStore()
	usersStore.AddUser("user1", "password")
	usersStore.AddUser("user2", "password")

	srv := NewDareServer(db, usersStore)

	authenticator := auth.NewJWTAutenticatorWithUsers(usersStore)
	mux := srv.CreateMux(auth.NewCasbinAuth(modelFile.Name(), policyFile.Name(), auth.Users{
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
	postData := map[string]string{"newKey": "newValue"}
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
		userStore: usersStore,
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

func TestHandlerCollectionGetById(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}
	srv.collectionManager.AddCollection("test-collection")
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	collection.Set("test-key", "test-value")

	req := httptest.NewRequest(http.MethodGet, "/test-collection/test-key", nil)
	w := httptest.NewRecorder()

	// Set PathValue as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")
	req.SetPathValue(KEY_PARAM, "test-key")

	// Execute
	srv.HandlerCollectionGetById(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", body["test-key"])
}

func TestHandlerCollectionSet(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}
	srv.collectionManager.AddCollection("test-collection")

	data := map[string]string{"test-key": "test-value"}
	jsonData, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/test-collection", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	// Set PathValue as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute
	srv.HandlerCollectionSet(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify that the key-value pair was stored
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	val := collection.Get("test-key")
	assert.Equal(t, "test-value", val)
}

func TestHandlerCollectionDelete(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}
	srv.collectionManager.AddCollection("test-collection")
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	collection.Set("test-key", "test-value")

	req := httptest.NewRequest(http.MethodDelete, "/test-collection/test-key", nil)
	w := httptest.NewRecorder()

	// Set PathValue as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")
	req.SetPathValue(KEY_PARAM, "test-key")

	// Execute
	srv.HandlerCollectionDelete(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the key was deleted
	val := collection.Get("test-key")
	assert.Equal(t, "", val)
}

func TestHandlerGetCollections(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}
	srv.collectionManager.AddCollection("collection1")
	srv.collectionManager.AddCollection("collection2")

	req := httptest.NewRequest(http.MethodGet, "/collections", nil)
	w := httptest.NewRecorder()

	// Execute
	srv.HandlerGetCollections(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var collections []string
	err := json.NewDecoder(resp.Body).Decode(&collections)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"collection1", "collection2"}, collections)
}

func TestHandlerCreateCollection(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}

	req := httptest.NewRequest(http.MethodPost, "/test-collection", nil)
	w := httptest.NewRecorder()

	// Set PathValue as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute
	srv.HandlerCreateCollection(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify the collection was created
	_, exists := srv.collectionManager.GetCollection("test-collection")
	assert.True(t, exists)
}

func TestHandlerDeleteCollection(t *testing.T) {
	// Setup
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}
	srv.collectionManager.AddCollection("test-collection")

	req := httptest.NewRequest(http.MethodDelete, "/test-collection", nil)
	w := httptest.NewRecorder()

	// Set PathValue as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute
	srv.HandlerDeleteCollection(w, req)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the collection was deleted
	_, exists := srv.collectionManager.GetCollection("test-collection")
	assert.False(t, exists)
}

// Test handler for successful pagination
func TestHandlerGetPaginatedItems_Success(t *testing.T) {
	collectionManager := database.NewCollectionManager()
	collectionManager.AddCollection(database.DEFAULT_COLLECTION)
	srv := &DareServer{
		collectionManager: collectionManager,
	}
	// Mock the default collection and add multiple items
	collection := srv.collectionManager.GetDefaultCollection()
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")
	collection.Set("key3", "value3")
	collection.Set("key4", "value4")

	// Simulate HTTP GET request for page 1 with 2 items per page
	req := httptest.NewRequest(http.MethodGet, "/items?page=1&pageSize=2", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	req.SetPathValue(COLLECTION_NAME_PARAM, "default")
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 2, "Expected 2 items in paginated response")

	assert.Equal(t, float64(1), responseBody["page"], "Expected page 1")
	assert.Equal(t, float64(2), responseBody["pageSize"], "Expected pageSize 2")
}

// Test handler for a different page
func TestHandlerGetPaginatedItems_Page2(t *testing.T) {
	collectionManager := database.NewCollectionManager()
	collectionManager.AddCollection(database.DEFAULT_COLLECTION)
	srv := &DareServer{
		collectionManager: collectionManager,
	}
	// Mock the default collection and add multiple items
	collection := srv.collectionManager.GetDefaultCollection()
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")
	collection.Set("key3", "value3")
	collection.Set("key4", "value4")

	// Simulate HTTP GET request for page 2 with 2 items per page
	req := httptest.NewRequest(http.MethodGet, "/items?page=2&pageSize=2", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	req.SetPathValue(COLLECTION_NAME_PARAM, "default")
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 2, "Expected 2 items in paginated response")

	assert.Equal(t, float64(2), responseBody["page"], "Expected page 2")
	assert.Equal(t, float64(2), responseBody["pageSize"], "Expected pageSize 2")
}

// Test handler when no items are available for the given page
func TestHandlerGetPaginatedItems_EmptyPage(t *testing.T) {
	collectionManager := database.NewCollectionManager()
	collectionManager.AddCollection(database.DEFAULT_COLLECTION)
	srv := &DareServer{
		collectionManager: collectionManager,
	}

	// Mock the default collection and add only a few items
	collection := srv.collectionManager.GetDefaultCollection()
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")

	// Simulate HTTP GET request for page 2 with 5 items per page (no data should exist for page 2)
	req := httptest.NewRequest(http.MethodGet, "/items?page=2&pageSize=5", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	req.SetPathValue(COLLECTION_NAME_PARAM, "default")
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 0, "Expected 0 items for empty page")

	assert.Equal(t, float64(2), responseBody["page"], "Expected page 2")
	assert.Equal(t, float64(5), responseBody["pageSize"], "Expected pageSize 5")
}

// Test handler for invalid method
func TestHandlerGetPaginatedItems_InvalidMethod(t *testing.T) {
	collectionManager := database.NewCollectionManager()
	collectionManager.AddCollection(database.DEFAULT_COLLECTION)
	srv := &DareServer{
		collectionManager: collectionManager,
	}

	// Simulate HTTP POST request (invalid method)
	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	req.SetPathValue(COLLECTION_NAME_PARAM, "default")
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Expected status 405 Method Not Allowed")
}

func TestHandlerGetPaginatedCollectionItems_Success(t *testing.T) {
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}

	// Mock a collection with multiple items
	srv.collectionManager.AddCollection("test-collection")
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")
	collection.Set("key3", "value3")
	collection.Set("key4", "value4")

	// Simulate HTTP GET request for page 1 with 2 items per page
	req := httptest.NewRequest(http.MethodGet, "/test-collection?page=1&pageSize=2", nil)
	w := httptest.NewRecorder()

	// Set path parameter as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute the handler
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 2, "Expected 2 items in paginated response")

	assert.Equal(t, float64(1), responseBody["page"], "Expected page 1")
	assert.Equal(t, float64(2), responseBody["pageSize"], "Expected pageSize 2")
}

// Test successful pagination for page 2 of a collection
func TestHandlerGetPaginatedCollectionItems_Page2(t *testing.T) {
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}

	// Mock a collection with multiple items
	srv.collectionManager.AddCollection("test-collection")
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")
	collection.Set("key3", "value3")
	collection.Set("key4", "value4")

	// Simulate HTTP GET request for page 2 with 2 items per page
	req := httptest.NewRequest(http.MethodGet, "/test-collection?page=2&pageSize=2", nil)
	w := httptest.NewRecorder()

	// Set path parameter as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute the handler
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 2, "Expected 2 items in paginated response")

	assert.Equal(t, float64(2), responseBody["page"], "Expected page 2")
	assert.Equal(t, float64(2), responseBody["pageSize"], "Expected pageSize 2")
}

// Test handler when no items are available for the given page
func TestHandlerGetPaginatedCollectionItems_EmptyPage(t *testing.T) {
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}

	// Mock a collection with 2 items
	srv.collectionManager.AddCollection("test-collection")
	collection, _ := srv.collectionManager.GetCollection("test-collection")
	collection.Set("key1", "value1")
	collection.Set("key2", "value2")

	// Simulate HTTP GET request for page 2 with 5 items per page (no data should exist for page 2)
	req := httptest.NewRequest(http.MethodGet, "/test-collection?page=2&pageSize=5", nil)
	w := httptest.NewRecorder()

	// Set path parameter as if from a router
	req.SetPathValue(COLLECTION_NAME_PARAM, "test-collection")

	// Execute the handler
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	items := responseBody["items"].([]interface{})
	assert.Len(t, items, 0, "Expected 0 items for empty page")

	assert.Equal(t, float64(2), responseBody["page"], "Expected page 2")
	assert.Equal(t, float64(5), responseBody["pageSize"], "Expected pageSize 5")
}

// Test handler for invalid method
func TestHandlerGetPaginatedCollectionItems_InvalidMethod(t *testing.T) {
	srv := &DareServer{
		collectionManager: database.NewCollectionManager(),
	}

	// Simulate HTTP POST request (invalid method)
	req := httptest.NewRequest(http.MethodPost, "/test-collection", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	srv.HandlerGetPaginatedCollectionItems(w, req)

	// Assert the response
	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Expected status 405 Method Not Allowed")
}
