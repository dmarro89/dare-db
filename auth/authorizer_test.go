package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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
g, user2, role2`

func TestNewCasbinAuth(t *testing.T) {
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

	require.NotNil(t, casbinAuth)
	require.NotNil(t, casbinAuth.enforcer)
	require.NotNil(t, casbinAuth.logger)
}

func TestHasPermission(t *testing.T) {
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

	// Test with unknown user
	ok := casbinAuth.HasPermission("unknown", "GET", "dare-db")
	require.False(t, ok)
}
