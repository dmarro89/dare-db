package auth

import (
	"log"
	"net/http"

	"github.com/casbin/casbin"
)

type Authorizer interface {
	HasPermission(userID, action, asset string) bool
}

type users struct {
	user  string
	roles []string
}

type CasbinAuth struct {
	users    map[string]users
	enforcer *casbin.Enforcer
}

func (a *CasbinAuth) HasPermission(userID, action, asset string) bool {
	user, ok := a.users[userID]
	if !ok {
		log.Print("Unknown user:", userID)
		return false
	}

	for _, role := range user.roles {
		if a.enforcer.Enforce(role, asset, action) {
			return true
		}
	}

	return false
}

type Middleware struct {
	authorizer Authorizer
}

func NewCasbinMiddleware() *Middleware {
	return &Middleware{
		authorizer: &CasbinAuth{},
	}
}

func (middleware *Middleware) HandleFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()
		asset := r.PathValue("asset")
		if !ok || !middleware.authorizer.HasPermission(username, r.Method, asset) {
			log.Printf("User '%s' is not allowed to '%s' resource '%s'", username, r.Method, asset)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
