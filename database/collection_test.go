package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCollectionManager(t *testing.T) {
	cm := NewCollectionManager()

	// Assert that cm is not nil
	assert.NotNil(t, cm, "Expected NewCollectionManager to return a non-nil value")

	// Assert that the collections map is initially empty
	assert.Equal(t, 0, len(cm.collections), "Expected initial collections map to be empty")
}

func TestAddCollection(t *testing.T) {
	cm := NewCollectionManager()
	cm.AddCollection("test-collection")

	// Assert that the collections map contains exactly one collection
	assert.Equal(t, 1, len(cm.collections), "Expected collections map to have 1 collection")

	// Assert that 'test-collection' was added to the collections map
	_, exists := cm.collections["test-collection"]
	assert.True(t, exists, "Expected collection 'test-collection' to be added")
}

func TestGetCollection(t *testing.T) {
	cm := NewCollectionManager()
	cm.AddCollection("test-collection")

	// Assert that the collection exists and is not nil
	db, exists := cm.GetCollection("test-collection")
	assert.True(t, exists, "Expected collection 'test-collection' to exist")
	assert.NotNil(t, db, "Expected non-nil database for collection 'test-collection'")
}

func TestGetCollectionNotExist(t *testing.T) {
	cm := NewCollectionManager()

	// Assert that a nonexistent collection returns false
	_, exists := cm.GetCollection("nonexistent-collection")
	assert.False(t, exists, "Expected nonexistent collection to return false")
}

func TestGetDefaultCollection(t *testing.T) {
	cm := NewCollectionManager()
	cm.AddCollection(DEFAULT_COLLECTION)

	// Assert that the default collection is not nil
	db := cm.GetDefaultCollection()
	assert.NotNil(t, db, "Expected non-nil database for default collection")
}

func TestGetCollectionNames(t *testing.T) {
	cm := NewCollectionManager()
	cm.AddCollection("collection1")
	cm.AddCollection("collection2")

	// Retrieve collection names and assert they match the expected names
	names := cm.GetCollectionNames()
	expectedNames := []string{"collection1", "collection2"}

	assert.ElementsMatch(t, expectedNames, names, "Expected all collection names to be returned")
}

func TestRemoveCollection(t *testing.T) {
	cm := NewCollectionManager()
	cm.AddCollection("collection-to-remove")

	// Remove the collection and assert the collections map is empty
	cm.RemoveCollection("collection-to-remove")
	assert.Equal(t, 0, len(cm.collections), "Expected collections map to be empty")

	// Assert that the collection no longer exists
	_, exists := cm.collections["collection-to-remove"]
	assert.False(t, exists, "Expected collection 'collection-to-remove' to be removed")
}
