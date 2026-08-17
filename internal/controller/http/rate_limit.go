package http

import (
	"errors"
	"net/http"
	"strconv"

	appmw "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

const defaultRateLimitPerMinute = 100

var errRateLimited = errors.New("too many requests")

func userRateLimitKey(r *http.Request) (string, error) {
	if userID, ok := appmw.UserID(r.Context()); ok {
		return "user:" + strconv.FormatInt(userID, 10), nil
	}

	return clientIPKey(r)
}

func clientIPKey(r *http.Request) (string, error) {
	return "ip:" + httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}

func rateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	render.Error(r.Context(), w, errRateLimited, http.StatusTooManyRequests, "rate limit exceeded")
}
