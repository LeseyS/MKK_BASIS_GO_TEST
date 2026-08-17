package v1

import (
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) TeamStats(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	stats, err := h.usecase.TeamStats(r.Context(), userID)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "failed to collect team stats")
		return
	}

	render.JSON(w, map[string]any{"teams": stats}, http.StatusOK)
}
