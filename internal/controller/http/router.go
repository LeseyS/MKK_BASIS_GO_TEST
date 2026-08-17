package http

import (
	"net/http"
	"time"

	ver1 "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/controller/http/v1"
	appmw "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/usecase"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/logger"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(r *chi.Mux, uc *usecase.UseCase, m *metrics.HTTPServer, jwtIssuer *jwtutil.Issuer) {
	v1 := ver1.New(uc)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api", func(r chi.Router) {
		r.Use(logger.Middleware)
		r.Use(metrics.NewMiddleware(m))

		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(httprate.LimitBy(100, time.Minute, func(r *http.Request) (string, error) {
					return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
				}))

				r.Post("/register", v1.CreateUser)
				r.Post("/login", v1.UserLogin)
			})

			r.Group(func(r chi.Router) {
				r.Use(appmw.Auth(jwtIssuer))

				r.Post("/teams", v1.CreateTeam)
				r.Get("/teams", v1.TeamListForUser)
				r.Get("/teams/stats", v1.TeamStats)
				r.Get("/teams/top-creators", v1.TeamTopCreators)
				r.Post("/teams/{id}/invite", v1.InviteUser)

				r.Post("/tasks", v1.CreateTask)
				r.Get("/tasks", v1.ListTasks)
				r.Get("/tasks/invalid-assignees", v1.TasksInvalidAssignee)
				r.Put("/tasks/{id}", v1.UpdateTask)
				r.Get("/tasks/{id}/history", v1.TaskHistory)
			})
		})
	})
}
