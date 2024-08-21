package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStore_AddUser(t *testing.T) {
	store := NewUserStore()

	// Test adding a new user
	err := store.AddUser("user1", "password1")
	assert.NoError(t, err, "Expected no error when adding a new user")

	// Test adding a duplicate user
	err = store.AddUser("user1", "password2")
	assert.Error(t, err, "Expected error when adding a duplicate user")
	assert.Equal(t, "user already exists", err.Error())
}

func TestUserStore_DeleteUser(t *testing.T) {
	store := NewUserStore()

	// Add a user to delete
	store.AddUser("user1", "password1")

	// Test deleting an existing user
	err := store.DeleteUser("user1")
	assert.NoError(t, err, "Expected no error when deleting an existing user")

	// Test deleting a non-existing user
	err = store.DeleteUser("user2")
	assert.Error(t, err, "Expected error when deleting a non-existing user")
	assert.Equal(t, "user does not exist", err.Error())
}

func TestUserStore_UpdatePassword(t *testing.T) {
	store := NewUserStore()

	// Add a user to update
	store.AddUser("user1", "password1")

	// Test updating the password for an existing user
	err := store.UpdatePassword("user1", "newpassword")
	assert.NoError(t, err, "Expected no error when updating password for an existing user")

	// Verify the password was updated
	assert.True(t, store.ValidateCredentials("user1", "newpassword"), "Expected the updated password to be valid")

	// Test updating the password for a non-existing user
	err = store.UpdatePassword("user2", "newpassword")
	assert.Error(t, err, "Expected error when updating password for a non-existing user")
	assert.Equal(t, "user does not exist", err.Error())
}

func TestUserStore_ValidateCredentials(t *testing.T) {
	store := NewUserStore()

	// Add a user to validate
	store.AddUser("user1", "password1")

	// Test validating correct credentials
	valid := store.ValidateCredentials("user1", "password1")
	assert.True(t, valid, "Expected credentials to be valid")

	// Test validating incorrect password
	valid = store.ValidateCredentials("user1", "wrongpassword")
	assert.False(t, valid, "Expected credentials to be invalid")

	// Test validating non-existing user
	valid = store.ValidateCredentials("user2", "password1")
	assert.False(t, valid, "Expected credentials to be invalid for non-existing user")
}

func TestUserStore_SaveToken(t *testing.T) {
	store := NewUserStore()

	// Save a token for a user
	store.SaveToken("user1", "token123")

	// Test that the token was saved correctly
	assert.True(t, store.ValidateToken("user1", "token123"), "Expected the token to be valid")
}

func TestUserStore_ValidateToken(t *testing.T) {
	store := NewUserStore()

	// Save a token for a user
	store.SaveToken("user1", "token123")

	// Test validating correct token
	valid := store.ValidateToken("user1", "token123")
	assert.True(t, valid, "Expected the token to be valid")

	// Test validating incorrect token
	valid = store.ValidateToken("user1", "wrongtoken")
	assert.False(t, valid, "Expected the token to be invalid")

	// Test validating token for non-existing user
	valid = store.ValidateToken("user2", "token123")
	assert.False(t, valid, "Expected the token to be invalid for non-existing user")
}
