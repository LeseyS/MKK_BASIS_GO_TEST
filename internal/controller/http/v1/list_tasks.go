package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
)

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	q := r.URL.Query()

	teamID, err := strconv.ParseInt(q.Get("team_id"), 10, 64)
	if err != nil || teamID == 0 {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "team_id query parameter is required")
		return
	}

	var assigneeID int64
	if v := r.URL.Query().Get("assignee_id"); v != "" {
		assigneeID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid assignee_id")
			return
		}
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	tasks, err := h.usecase.ListTasks(r.Context(), dto.ListTasks{
		TeamID:     teamID,
		UserID:     userID,
		Status:     q.Get("status"),
		AssigneeID: assigneeID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeTaskServiceError(w, r, err)
		return
	}

	render.JSON(w, map[string]any{
		"tasks":  tasks.Tasks,
		"total":  tasks.Total,
		"limit":  tasks.Limit,
		"offset": tasks.Offset,
	}, http.StatusOK)
}

func writeTaskServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, apperr.ErrForbidden):
		render.Error(r.Context(), w, err, http.StatusForbidden, "you are not a member of this team")
	case errors.Is(err, apperr.ErrInvalidAssignee):
		render.Error(r.Context(), w, err, http.StatusUnprocessableEntity, "assignee must be a member of the task's team")
	case errors.Is(err, apperr.ErrNotFound):
		render.Error(r.Context(), w, err, http.StatusNotFound, "task not found")
	case errors.Is(err, apperr.ErrValidation):
		render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid request body")
	default:
		render.Error(r.Context(), w, err, http.StatusInternalServerError, "internal error")
	}
}
