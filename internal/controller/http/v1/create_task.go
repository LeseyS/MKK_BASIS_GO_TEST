package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	var req dto.CreateTaskIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TeamID == 0 || req.Title == "" {
		render.Error(r.Context(), w, errors.New("invalid request body"), http.StatusBadRequest, "invalid request body")
		return
	}

	req.ActorID = userID

	created, err := h.usecase.CreateTask(r.Context(), req)
	if err != nil {
		writeTaskServiceError(w, r, err)
		return
	}

	render.JSON(w, created, http.StatusCreated)
}
