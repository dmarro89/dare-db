package auth

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/dmarro89/dare-db/logger"
)

type Authorizer interface {
	HasPermission(userID, action, asset string) bool
}

type User struct {
	User  string
	Roles []string
}

type CasbinAuth struct {
	users    map[string]User
	enforcer *casbin.Enforcer
	logger   logger.Logger
}

func NewCasbinAuth(modelPath, policyPath string, users map[string]User) *CasbinAuth {
	if users == nil {
		users = map[string]User{"guest": {User: "guest", Roles: []string{"guest"}}}
	}
	enforcer := casbin.NewEnforcer(modelPath, policyPath)
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
			a.logger.Info(fmt.Printf("User '%s' is allowed to '%s' resource '%s'", userID, action, asset))
			return true
		}
	}

	a.logger.Info(fmt.Printf(`User %s is not allowed to %s resource %s`, userID, action, asset))
	return false
}

type Middleware struct {
	authorizer Authorizer
	logger     logger.Logger
}

func NewCasbinMiddleware(casbinAuth *CasbinAuth) *Middleware {
	return &Middleware{
		authorizer: casbinAuth,
		logger:     logger.NewDareLogger(),
	}
}

func (middleware *Middleware) HandleFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()
		if !ok {
			middleware.logger.Info("Missing or invalid credentials")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		asset := r.PathValue("asset")
		if asset == "" {
			asset = "dare-db"
		}
		if !middleware.authorizer.HasPermission(username, r.Method, asset) {
			middleware.logger.Info(fmt.Printf(`User %s is not allowed to %s resource %s`, username, r.Method, asset))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r)
	})
}
