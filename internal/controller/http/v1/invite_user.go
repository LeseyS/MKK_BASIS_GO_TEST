package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) InviteUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	teamID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dto.InviteUserIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if err := validate.V.Struct(req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "validation error")

		return
	}

	req.UserID = userID
	req.TeamID = teamID

	emailSent, err := h.usecase.InviteUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			render.Error(r.Context(), w, err, http.StatusForbidden, "only team owner/admin can invite members")
		case errors.Is(err, apperr.ErrAlreadyMember):
			render.Error(r.Context(), w, err, http.StatusConflict, "user is already a member of this team")
		default:
			render.Error(r.Context(), w, err, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	render.JSON(w, dto.InviteUserOut{
		Invited:   true,
		EmailSent: emailSent,
	}, http.StatusOK)
}
