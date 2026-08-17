package task_client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) CreateTask(ctx context.Context, in CreateTaskRequest) (CreateTaskResponse, error) {
	var out CreateTaskResponse

	err := c.do(ctx, http.MethodPost, apiPrefix+"/tasks", nil, in, &out)

	return out, err
}

func (c *Client) ListTasks(ctx context.Context, in ListTasksRequest) (ListTasksResponse, error) {
	var out ListTasksResponse

	query := url.Values{}
	query.Set("team_id", strconv.FormatInt(in.TeamID, 10))

	if in.Status != "" {
		query.Set("status", in.Status)
	}
	if in.AssigneeID != 0 {
		query.Set("assignee_id", strconv.FormatInt(in.AssigneeID, 10))
	}
	if in.Limit != 0 {
		query.Set("limit", strconv.Itoa(in.Limit))
	}
	if in.Offset != 0 {
		query.Set("offset", strconv.Itoa(in.Offset))
	}

	err := c.do(ctx, http.MethodGet, apiPrefix+"/tasks", query, nil, &out)

	return out, err
}

func (c *Client) UpdateTask(ctx context.Context, taskID int64, in UpdateTaskRequest) (UpdateTaskResponse, error) {
	var out UpdateTaskResponse

	path := apiPrefix + "/tasks/" + strconv.FormatInt(taskID, 10)
	err := c.do(ctx, http.MethodPut, path, nil, in, &out)

	return out, err
}

func (c *Client) TaskHistory(ctx context.Context, taskID int64, limit, offset int) (TaskHistoryResponse, error) {
	var out TaskHistoryResponse

	query := url.Values{}
	if limit != 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if offset != 0 {
		query.Set("offset", strconv.Itoa(offset))
	}

	path := apiPrefix + "/tasks/" + strconv.FormatInt(taskID, 10) + "/history"
	err := c.do(ctx, http.MethodGet, path, query, nil, &out)

	return out, err
}

func (c *Client) TasksInvalidAssignees(ctx context.Context) ([]TaskInvalidAssignee, error) {
	var out invalidAssigneesResponse

	err := c.do(ctx, http.MethodGet, apiPrefix+"/tasks/invalid-assignees", nil, nil, &out)

	return out.Tasks, err
}
