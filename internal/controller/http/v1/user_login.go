package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := validate.V.Struct(req); err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "validation error")
		return
	}

	userLogin, err := h.usecase.UserLogin(r.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			render.Error(r.Context(), w, err, http.StatusUnauthorized, "invalid email or password")
			return
		}

		render.Error(r.Context(), w, err, http.StatusInternalServerError, "failed to log in")
		return
	}

	render.JSON(w, dto.UserLoginOut{
		Token: userLogin.Token,
		User:  userLogin.User,
	}, http.StatusOK)
}
