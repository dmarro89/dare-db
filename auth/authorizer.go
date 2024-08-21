package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/casbin/casbin"
	"github.com/dmarro89/dare-db/logger"
)

const GUEST_USER = "guest"
const GUEST_ROLE = "guest"
const DEFAULT_USER = "admin"
const DEFAULT_ROLE = "admin"

type Authorizer interface {
	HasPermission(userID, action, asset string) bool
}

type User struct {
	Roles []string
}

type Users map[string]*User

type CasbinAuth struct {
	users    Users
	enforcer *casbin.Enforcer
	logger   logger.Logger
}

func NewCasbinAuth(modelPath, policyPath string, users Users) *CasbinAuth {
	if users == nil {
		users = Users{GUEST_USER: {Roles: []string{GUEST_ROLE}}}
	}
	enforcer, err := casbin.NewEnforcerSafe(modelPath, policyPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Casbin enforcer: %v", err))
	}
	return &CasbinAuth{
		users:    users,
		enforcer: enforcer,
		logger:   logger.NewDareLogger(),
	}
}

func (a *CasbinAuth) HasPermission(userID, action, asset string) bool {
	user, ok := a.users[userID]
	if !ok {
		a.logger.Error("Unknown user:", userID)
		return false
	}

	for _, role := range user.Roles {
		if a.enforcer.Enforce(role, asset, action) {
			a.logger.Info(fmt.Sprintf("User '%s' is allowed to '%s' resource '%s'", userID, action, asset))
			return true
		}
	}

	a.logger.Info(fmt.Sprintf("User %s is not allowed to %s resource %s", userID, action, asset))
	return false
}

func GetDefaultAuth() *CasbinAuth {
	dir, err := os.Getwd()
	if err != nil {
		panic("Failed to get current working directory: " + err.Error())
	}
	modelPath := filepath.Join(dir, "auth/rbac_model.conf")
	policyPath := filepath.Join(dir, "auth/rbac_policy.csv")

	return NewCasbinAuth(modelPath, policyPath, Users{
		DEFAULT_USER: {Roles: []string{DEFAULT_ROLE}},
	})
}
