package task_client

import (
	"context"
	"net/http"
	"strconv"
)

func (c *Client) CreateTeam(ctx context.Context, in CreateTeamRequest) (CreateTeamResponse, error) {
	var out CreateTeamResponse

	err := c.do(ctx, http.MethodPost, apiPrefix+"/teams", nil, in, &out)

	return out, err
}

func (c *Client) Teams(ctx context.Context) ([]TeamListItem, error) {
	var out []TeamListItem

	err := c.do(ctx, http.MethodGet, apiPrefix+"/teams", nil, nil, &out)

	return out, err
}

func (c *Client) InviteUser(ctx context.Context, teamID int64, in InviteRequest) (InviteResponse, error) {
	var out InviteResponse

	path := apiPrefix + "/teams/" + strconv.FormatInt(teamID, 10) + "/invite"
	err := c.do(ctx, http.MethodPost, path, nil, in, &out)

	return out, err
}

func (c *Client) TeamStats(ctx context.Context) ([]TeamStats, error) {
	var out teamStatsResponse

	err := c.do(ctx, http.MethodGet, apiPrefix+"/teams/stats", nil, nil, &out)

	return out.Teams, err
}

func (c *Client) TeamTopCreators(ctx context.Context) ([]TeamTopCreator, error) {
	var out topCreatorsResponse

	err := c.do(ctx, http.MethodGet, apiPrefix+"/teams/top-creators", nil, nil, &out)

	return out.TopCreators, err
}
