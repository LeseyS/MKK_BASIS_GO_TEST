package v1

import (
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) TasksInvalidAssignee(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	tasks, err := h.usecase.TasksInvalidAssignee(r.Context(), userID)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "failed to run integrity check")
		return
	}

	render.JSON(w, map[string]any{"tasks": tasks}, http.StatusOK)
}
