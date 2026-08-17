//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestOutsiderIsLockedOutOfTeam() {
	alice, _ := s.newUser("insider")
	bob, bobUser := s.newUser("stranger")

	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "внутренняя задача")

	s.Run("чтение списка задач", func() {
		_, err := bob.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
		s.ErrorIs(err, task_client.ErrForbidden)
	})

	s.Run("создание задачи", func() {
		_, err := bob.CreateTask(ctx, task_client.CreateTaskRequest{
			TeamID: teamID, Title: "чужая задача",
		})
		s.ErrorIs(err, task_client.ErrForbidden)
	})

	s.Run("правка задачи", func() {
		status := "done"
		_, err := bob.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &status})
		s.ErrorIs(err, task_client.ErrForbidden)
	})

	s.Run("чтение истории", func() {
		_, err := bob.TaskHistory(ctx, task.ID, 0, 0)
		s.ErrorIs(err, task_client.ErrForbidden)
	})

	s.Run("приглашение в команду", func() {
		_, err := bob.InviteUser(ctx, teamID, task_client.InviteRequest{Email: bobUser.Email})
		s.ErrorIs(err, task_client.ErrForbidden)
	})

	s.Run("задача осталась нетронутой", func() {
		out, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
		s.NoError(err)
		s.Equal(int64(1), out.Total)
		s.Equal("todo", out.Tasks[0].Status)
	})
}

func (s *Suite) TestMemberHasFullTaskAccess() {
	alice, _ := s.newUser("teamlead")
	bob, bobUser := s.newUser("teammate")

	teamID := s.newTeam(alice, "Alpha")
	s.invite(alice, teamID, bobUser.Email)

	task := s.newTask(alice, teamID, "общая задача")

	list, err := bob.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(1), list.Total)

	status := "in_progress"
	updated, err := bob.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &status})
	s.NoError(err)
	s.Equal("in_progress", updated.Task.Status)

	history, err := bob.TaskHistory(ctx, task.ID, 0, 0)
	s.NoError(err)
	s.Equal(int64(1), history.Total)
	s.Equal(bobUser.ID, history.History[0].ChangedBy)
}

func (s *Suite) TestAnalyticsAreScopedToOwnTeams() {
	alice, _ := s.newUser("analyticsowner")
	bob, _ := s.newUser("analyticsstranger")

	teamID := s.newTeam(alice, "Alpha")
	s.newTask(alice, teamID, "задача")

	stats, err := bob.TeamStats(ctx)
	s.NoError(err)
	for _, row := range stats {
		s.NotEqual(teamID, row.TeamID, "чужая команда не должна попадать в статистику")
	}

	top, err := bob.TeamTopCreators(ctx)
	s.NoError(err)
	for _, row := range top {
		s.NotEqual(teamID, row.TeamID)
	}

	invalid, err := bob.TasksInvalidAssignees(ctx)
	s.NoError(err)
	for _, row := range invalid {
		s.NotEqual(teamID, row.TeamID)
	}
}
