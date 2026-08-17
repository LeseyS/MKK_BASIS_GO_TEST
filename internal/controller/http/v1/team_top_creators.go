package v1

import (
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) TeamTopCreators(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	top, err := h.usecase.TeamTopCreators(r.Context(), userID)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "failed to collect top creators")
		return
	}

	render.JSON(w, map[string]any{"top_creators": top}, http.StatusOK)
}
