package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/passwordhash"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateUserIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error")

		return
	}

	if err := validate.V.Struct(req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "validation error")

		return
	}

	user, err := h.usecase.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrEmailTaken):
			render.Error(ctx, w, err, http.StatusConflict, "email already registered")
		case errors.Is(err, apperr.ErrValidation), errors.Is(err, passwordhash.ErrPasswordTooLong):
			render.Error(ctx, w, err, http.StatusBadRequest, "validation error")
		default:
			render.Error(ctx, w, err, http.StatusInternalServerError, "create user error")
		}

		return
	}

	render.JSON(w, user, http.StatusCreated)
}
