package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dmarro89/dare-db/logger"
)

type Middleware interface {
	HandleFunc(next http.HandlerFunc) http.HandlerFunc
}
type DareMiddleware struct {
	authorizer    Authorizer
	authenticator Authenticator
	logger        logger.Logger
}

func NewCasbinMiddleware(casbinAuth Authorizer, authenticator Authenticator) Middleware {
	return &DareMiddleware{
		authorizer:    casbinAuth,
		authenticator: authenticator,
		logger:        logger.NewDareLogger(),
	}
}

func (middleware *DareMiddleware) HandleFunc(next http.HandlerFunc) http.HandlerFunc {
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
			return
		}

		next(w, r)
	})
}

func (middleware *DareMiddleware) extractAssetFromPath(path string) string {
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
