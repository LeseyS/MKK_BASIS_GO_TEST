//go:build integration

package test

import (
	"fmt"
	"sync/atomic"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"
)

var userCounter atomic.Int64

const testPassword = "password123"

func (s *Suite) newUser(name string) (*task_client.Client, task_client.User) {
	email := fmt.Sprintf("%s-%d@example.com", name, userCounter.Add(1))

	client, user, err := s.api.RegisterAndLogin(ctx, name, email, testPassword)
	s.NoError(err)

	return client, user
}

func (s *Suite) newTeam(owner *task_client.Client, name string) int64 {
	out, err := owner.CreateTeam(ctx, task_client.CreateTeamRequest{
		Name: fmt.Sprintf("%s-%d", name, userCounter.Add(1)),
	})
	s.NoError(err)

	return out.Team.ID
}

func (s *Suite) invite(owner *task_client.Client, teamID int64, email string) {
	out, err := owner.InviteUser(ctx, teamID, task_client.InviteRequest{Email: email})
	s.NoError(err)
	s.True(out.Invited)
}

func (s *Suite) newTask(client *task_client.Client, teamID int64, title string) task_client.Task {
	out, err := client.CreateTask(ctx, task_client.CreateTaskRequest{
		TeamID: teamID,
		Title:  title,
	})
	s.NoError(err)

	return out.Task
}
