package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/dmarro89/dare-db/logger"
)

const GUEST_USER = "guest"
const GUEST_ROLE = "guest"

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

type Middleware struct {
	authorizer    Authorizer
	authenticator Authenticator
	logger        logger.Logger
}

func NewCasbinMiddleware(casbinAuth Authorizer, authenticator Authenticator) *Middleware {
	return &Middleware{
		authorizer:    casbinAuth,
		authenticator: authenticator,
		logger:        logger.NewDareLogger(),
	}
}

func (middleware *Middleware) HandleFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			middleware.logger.Info("Missing authorization token")
			http.Error(w, "Unauthorized: missing authorization token", http.StatusUnauthorized)
			return
		}

		username, err := middleware.authenticator.VerifyToken(tokenStr)
		if err != nil {
			middleware.logger.Error(fmt.Sprintf("Invalid authorization token: %v", err))
			http.Error(w, "Unauthorized: invalid authorization token", http.StatusUnauthorized)
		}

		asset := middleware.extractAssetFromPath(r.URL.Path)

		middleware.logger.Info(fmt.Sprintf("User %s is requesting %s resource %s", username, r.Method, asset))
		if !middleware.authorizer.HasPermission(username, r.Method, asset) {
			middleware.logger.Info(fmt.Sprintf("User %s is not allowed to %s resource %s", username, r.Method, asset))
			http.Error(w, "Forbidden: you do not have permission to access this resource", http.StatusForbidden)
			return
		}

		next(w, r)
	})
}

func (middleware *Middleware) extractAssetFromPath(path string) string {
	if strings.HasPrefix(path, "/get/") {
		return strings.TrimPrefix(path, "/get/")
	}
	if strings.HasPrefix(path, "/delete/") {
		return strings.TrimPrefix(path, "/delete/")
	}
	if path == "/set" {
		return "set"
	}
	return "dare-db"
}
