package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req dto.UpdateTaskIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.V.Struct(req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "validation error")
		return
	}

	req.ActorID = userID
	req.TaskID = taskID

	upd, err := h.usecase.UpdateTask(r.Context(), req)
	if err != nil {
		writeTaskServiceError(w, r, err)
		return
	}

	render.JSON(w, upd, http.StatusOK)

}
