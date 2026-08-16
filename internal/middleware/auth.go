package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

type ctxKey string

const userIDKey ctxKey = "userID"
const userEmailKey ctxKey = "userEmail"

func Auth(issuer *jwtutil.Issuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				render.JSON(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				render.JSON(w, "Authorization header must be 'Bearer <token>'", http.StatusUnauthorized)
				return
			}

			claims, err := issuer.Parse(parts[1])
			if err != nil {
				render.JSON(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userEmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

func UserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}
