package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmarro89/dare-db/database"
	"gotest.tools/assert"
)

func TestServer_SetAndGet(t *testing.T) {
	db := database.NewDatabase()
	srv := NewServer(db)

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
	srv.HandlerGet(getEmptyResponse, getEmptyRequest)
	assert.Equal(t, http.StatusBadRequest, getEmptyResponse.Code, `url param query "key" cannot be empty`)

	getMissingKeyRequest, _ := http.NewRequest("GET", "/get?key=missingKey", nil)
	getMissingKeyResponse := httptest.NewRecorder()
	srv.HandlerGet(getMissingKeyResponse, getMissingKeyRequest)
	assert.Equal(t, http.StatusNotFound, getMissingKeyResponse.Code, `Key "%s" not found`, "missingKey")

	getRequest, _ := http.NewRequest("GET", "/get?key=testKey", nil)
	getResponse := httptest.NewRecorder()

	srv.HandlerGet(getResponse, getRequest)

	assert.Equal(t, http.StatusOK, getResponse.Code, "Expected status %d, got %d", http.StatusOK, getResponse.Code)

	var getResult map[string]interface{}
	err := json.Unmarshal(getResponse.Body.Bytes(), &getResult)
	assert.NilError(t, err, "Error decoding JSON response")

	for key, value := range getResult {
		assert.Equal(t, key, "testKey", "Unexpected response body: %v", getResult)
		assert.Equal(t, value, "testValue", "Unexpected response body: %v", getResult)
	}
}

func TestServer_SetAndDelete(t *testing.T) {
	db := database.NewDatabase()
	srv := NewServer(db)

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

	deleteRequest, _ := http.NewRequest("DELETE", "/delete?key=testKey", nil)
	deleteResponse := httptest.NewRecorder()

	srv.HandlerDelete(deleteResponse, deleteRequest)

	if deleteResponse.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, deleteResponse.Code)
	}

	getRequest, _ := http.NewRequest("GET", "/get?key=testKey", nil)
	getResponse := httptest.NewRecorder()

	srv.HandlerGet(getResponse, getRequest)

	if getResponse.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, getResponse.Code)
	}
}
