// database_test.go

package database

import (
	"testing"
)

func TestDatabase_SetAndGet(t *testing.T) {
	db := NewDatabase()

	key := "testKey"
	value := "testValue"

	err := db.Set(key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	result := db.Get(key)
	if result != value {
		t.Errorf("Expected %v, got %v", value, result)
	}
}

func TestDatabase_SetAndGetConcurrently(t *testing.T) {
	db := NewDatabase()

	key := "testKey"
	value := "testValue"

	go func() {
		err := db.Set(key, value)
		if err != nil {
			t.Errorf("Error setting value: %v", err)
		}
	}()

	go func() {
		result := db.Get(key)
		if result != value {
			t.Errorf("Expected %v, got %v", value, result)
		}
	}()

	t.Logf("Waiting for goroutines to finish...")
}

func TestDatabase_Delete(t *testing.T) {
	db := NewDatabase()

	key := "testKey"
	value := "testValue"

	err := db.Set(key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	err = db.Delete(key)
	if err != nil {
		t.Errorf("Error deleting key: %v", err)
	}

	result := db.Get(key)
	if result != nil {
		t.Errorf("Expected nil after deletion, got %v", result)
	}
}
