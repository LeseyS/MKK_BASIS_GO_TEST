package v1

import (
	"net/http"
	"strconv"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/middleware"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/render"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) TaskHistory(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid task ID")
		return
	}

	q := r.URL.Query()

	var limit int
	if v := q.Get("limit"); v != "" {
		limit, err = strconv.Atoi(v)
		if err != nil {
			render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid limit")
			return
		}
	}

	var offset int
	if v := q.Get("offset"); v != "" {
		offset, err = strconv.Atoi(v)
		if err != nil {
			render.Error(r.Context(), w, err, http.StatusBadRequest, "invalid offset")
			return
		}
	}

	history, err := h.usecase.TaskHistory(r.Context(), dto.TaskHistoryIn{
		TaskID:  taskID,
		ActorID: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		writeTaskServiceError(w, r, err)
		return
	}

	render.JSON(w, map[string]any{
		"history": history.Entries,
		"total":   history.Total,
		"limit":   history.Limit,
		"offset":  history.Offset,
	}, http.StatusOK)
}
