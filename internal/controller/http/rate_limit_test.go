package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appmw "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/stretchr/testify/require"
)

const testLimit = 3

func newLimitedRouter(issuer *jwtutil.Issuer) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.ClientIPFromRemoteAddr)

	r.Group(func(r chi.Router) {
		r.Use(appmw.Auth(issuer))
		r.Use(httprate.Limit(testLimit, time.Minute,
			httprate.WithKeyFuncs(userRateLimitKey),
			httprate.WithLimitHandler(rateLimitExceeded),
		))

		r.Get("/protected", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return r
}

func call(t *testing.T, h http.Handler, token string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	return rec
}

func TestRateLimit_PerUserBuckets(t *testing.T) {
	issuer := jwtutil.NewIssuer("secret", time.Hour)
	router := newLimitedRouter(issuer)

	alice, err := issuer.Generate(1, "alice@example.com")
	require.NoError(t, err)

	bob, err := issuer.Generate(2, "bob@example.com")
	require.NoError(t, err)

	for i := 0; i < testLimit; i++ {
		require.Equal(t, http.StatusOK, call(t, router, alice).Code, "запрос %d в пределах лимита", i+1)
	}

	require.Equal(t, http.StatusTooManyRequests, call(t, router, alice).Code,
		"превышение лимита должно отдавать 429")

	require.Equal(t, http.StatusOK, call(t, router, bob).Code,
		"у второго пользователя своя корзина, общий IP роли не играет")
}

func TestRateLimit_ResponseIsJSON(t *testing.T) {
	issuer := jwtutil.NewIssuer("secret", time.Hour)
	router := newLimitedRouter(issuer)

	token, err := issuer.Generate(1, "alice@example.com")
	require.NoError(t, err)

	for i := 0; i < testLimit; i++ {
		call(t, router, token)
	}

	rec := call(t, router, token)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Contains(t, body.Error, "rate limit exceeded")
}

func TestRateLimit_UnauthenticatedFallsBackToIP(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(httprate.Limit(testLimit, time.Minute,
		httprate.WithKeyFuncs(userRateLimitKey),
		httprate.WithLimitHandler(rateLimitExceeded),
	))
	r.Get("/public", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	send := func(ip string) int {
		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.RemoteAddr = ip + ":1234"

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		return rec.Code
	}

	for i := 0; i < testLimit; i++ {
		require.Equal(t, http.StatusOK, send("192.0.2.1"))
	}

	require.Equal(t, http.StatusTooManyRequests, send("192.0.2.1"),
		"без токена ключом становится IP")
	require.Equal(t, http.StatusOK, send("192.0.2.9"),
		"другой IP не должен попадать под чужой лимит")
}
