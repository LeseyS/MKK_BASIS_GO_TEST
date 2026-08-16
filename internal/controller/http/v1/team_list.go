package v1

import (
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) TeamListForUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	teams, err := h.usecase.TeamListForUser(r.Context(), userID)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "failed to list teams")
		return
	}

	render.JSON(w, teams, http.StatusOK)
}
