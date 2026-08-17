package domain

import (
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask_DefaultStatus(t *testing.T) {
	task, err := NewTask(1, "title", "description", "", 1, nil)

	require.NoError(t, err)
	assert.Equal(t, StatusTodo, task.Status, "пустой статус превращается в todo")
}

func TestNewTask_StatusValidation(t *testing.T) {
	cases := []struct {
		status  TaskStatus
		wantErr bool
	}{
		{StatusTodo, false},
		{StatusInProgress, false},
		{StatusDone, false},
		{"bogus", true},
		{"DONE", true},
	}

	for _, c := range cases {
		t.Run(string(c.status), func(t *testing.T) {
			_, err := NewTask(1, "title", "", c.status, 1, nil)

			if c.wantErr {
				require.ErrorIs(t, err, apperr.ErrValidation)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewTask_TitleRequired(t *testing.T) {
	_, err := NewTask(1, "", "description", StatusTodo, 1, nil)

	require.ErrorIs(t, err, apperr.ErrValidation)
}

func TestNewTask_KeepsFields(t *testing.T) {
	assignee := int64(5)

	task, err := NewTask(7, "title", "description", StatusInProgress, 3, &assignee)

	require.NoError(t, err)
	assert.Equal(t, int64(7), task.TeamID)
	assert.Equal(t, int64(3), task.CreatedBy)
	require.NotNil(t, task.AssigneeID)
	assert.Equal(t, assignee, *task.AssigneeID)
}
