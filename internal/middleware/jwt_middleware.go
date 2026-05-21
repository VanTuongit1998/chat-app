package middleware

import (
	"context"
	"net/http"
	"strings"

	"chat-app/internal/model"
	"chat-app/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func JwtAuthentication(jwtService *utils.JwtService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenString := ""

			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString == authHeader {
					utils.Error(w, http.StatusUnauthorized, "invalid authorization header")
					return
				}
			} else {
				tokenString = r.URL.Query().Get("token")
			}

			if tokenString == "" {
				utils.Error(w, http.StatusUnauthorized, "authorization token missing")
				return
			}

			user, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthenticatedUser(r *http.Request) *model.User {
	user, _ := r.Context().Value(UserContextKey).(*model.User)
	return user
}
