package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const modelPath = "./test/rbac_model.conf"
const policyPath = "./test/rbac_policy.csv"

func TestNewCasbinAuth(t *testing.T) {
	casbinAuth := NewCasbinAuth(modelPath, policyPath, Users{
		"user1": {Roles: []string{"role1"}},
		"user2": {Roles: []string{"role2"}},
	})

	require.NotNil(t, casbinAuth)
	require.NotNil(t, casbinAuth.enforcer)
	require.NotNil(t, casbinAuth.logger)
}

func TestHasPermission(t *testing.T) {
	modelPath := "./test/rbac_model.conf"
	policyPath := "./test/rbac_policy.csv"

	casbinAuth := NewCasbinAuth(modelPath, policyPath, Users{
		"user1": {Roles: []string{"role1"}},
		"user2": {Roles: []string{"role2"}},
	})

	// Test with user1 and GET on dare-db
	ok := casbinAuth.HasPermission("user1", "GET", "dare-db")
	require.True(t, ok)

	// Test with user2 and POST on dare-db
	ok = casbinAuth.HasPermission("user2", "POST", "dare-db")
	require.True(t, ok)

	// Test with user1 and POST on dare-db (should not have permission)
	ok = casbinAuth.HasPermission("user1", "POST", "dare-db")
	require.False(t, ok)

	// Test with user2 and GET on dare-db (should not have permission)
	ok = casbinAuth.HasPermission("user2", "GET", "dare-db")
	require.False(t, ok)
}

func TestUnknownUser(t *testing.T) {
	casbinAuth := NewCasbinAuth(modelPath, policyPath, Users{
		"user1": {Roles: []string{"role1"}},
		"user2": {Roles: []string{"role2"}},
	})

	// Test with unknown user
	ok := casbinAuth.HasPermission("unknown", "GET", "dare-db")
	require.False(t, ok)
}
