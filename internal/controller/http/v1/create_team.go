package v1

import (
	"encoding/json"
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
)

func (h *Handlers) CreateTeam(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	var req dto.CreateTeamIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "validation error")
		return
	}

	req.OwnerID = userID

	if err := validate.V.Struct(req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "validation error")
		return
	}

	team, err := h.usecase.CreateWithOwner(r.Context(), req)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "error creating team")
		return
	}

	render.JSON(w, team, http.StatusCreated)
}
