//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestTeamStats() {
	alice, _ := s.newUser("statsowner")
	_, bobUser := s.newUser("statsmember")
	_, carolUser := s.newUser("statsmember2")

	teamID := s.newTeam(alice, "Stats")
	s.invite(alice, teamID, bobUser.Email)
	s.invite(alice, teamID, carolUser.Email)

	first := s.newTask(alice, teamID, "закрытая 1")
	second := s.newTask(alice, teamID, "закрытая 2")
	s.newTask(alice, teamID, "открытая")

	done := "done"
	for _, id := range []int64{first.ID, second.ID} {
		_, err := alice.UpdateTask(ctx, id, task_client.UpdateTaskRequest{Status: &done})
		s.NoError(err)
	}

	stats, err := alice.TeamStats(ctx)
	s.NoError(err)

	var row *task_client.TeamStats
	for i := range stats {
		if stats[i].TeamID == teamID {
			row = &stats[i]
			break
		}
	}
	s.NotNil(row, "команда должна быть в статистике")

	s.Equal(int64(3), row.MembersCount, "владелец плюс два приглашённых")
	s.Equal(int64(2), row.DoneLast7Days)
}

func (s *Suite) TestTeamStatsCountsTeamWithoutTasks() {
	alice, _ := s.newUser("statsempty")
	teamID := s.newTeam(alice, "Empty")

	stats, err := alice.TeamStats(ctx)
	s.NoError(err)

	var found bool
	for _, row := range stats {
		if row.TeamID == teamID {
			found = true
			s.Equal(int64(1), row.MembersCount)
			s.Equal(int64(0), row.DoneLast7Days)
		}
	}
	s.True(found, "команда без задач всё равно должна попадать в выдачу")
}

func (s *Suite) TestTeamTopCreators() {
	alice, aliceUser := s.newUser("topowner")
	bob, bobUser := s.newUser("topsecond")
	carol, carolUser := s.newUser("topthird")
	dave, daveUser := s.newUser("topfourth")

	teamID := s.newTeam(alice, "Top")
	for _, email := range []string{bobUser.Email, carolUser.Email, daveUser.Email} {
		s.invite(alice, teamID, email)
	}

	creators := []struct {
		client *task_client.Client
		count  int
	}{
		{alice, 4},
		{bob, 3},
		{carol, 2},
		{dave, 1},
	}
	for _, c := range creators {
		for i := 0; i < c.count; i++ {
			s.newTask(c.client, teamID, "задача")
		}
	}

	top, err := alice.TeamTopCreators(ctx)
	s.NoError(err)

	var rows []task_client.TeamTopCreator
	for _, row := range top {
		if row.TeamID == teamID {
			rows = append(rows, row)
		}
	}

	s.Len(rows, 3, "в выдачу попадают только первые три создателя")

	s.Equal(1, rows[0].Position)
	s.Equal(aliceUser.ID, rows[0].UserID)
	s.Equal(int64(4), rows[0].TasksCreated)

	s.Equal(2, rows[1].Position)
	s.Equal(bobUser.ID, rows[1].UserID)
	s.Equal(int64(3), rows[1].TasksCreated)

	s.Equal(3, rows[2].Position)
	s.Equal(carolUser.ID, rows[2].UserID)
	s.Equal(int64(2), rows[2].TasksCreated)

	for _, row := range rows {
		s.NotEqual(daveUser.ID, row.UserID, "четвёртый создатель в топ-3 попадать не должен")
	}
}

func (s *Suite) TestTasksInvalidAssignee() {
	alice, _ := s.newUser("integrityowner")
	_, bobUser := s.newUser("integritymember")

	teamID := s.newTeam(alice, "Integrity")
	s.invite(alice, teamID, bobUser.Email)

	out, err := alice.CreateTask(ctx, task_client.CreateTaskRequest{
		TeamID: teamID, Title: "назначена на bob", AssigneeID: &bobUser.ID,
	})
	s.NoError(err)

	before, err := alice.TasksInvalidAssignees(ctx)
	s.NoError(err)
	for _, row := range before {
		s.NotEqual(out.Task.ID, row.TaskID, "пока bob в команде, нарушения нет")
	}

	_, err = s.db.ExecContext(ctx,
		"DELETE FROM team_members WHERE team_id = ? AND user_id = ?", teamID, bobUser.ID)
	s.NoError(err)

	after, err := alice.TasksInvalidAssignees(ctx)
	s.NoError(err)

	var found *task_client.TaskInvalidAssignee
	for i := range after {
		if after[i].TaskID == out.Task.ID {
			found = &after[i]
			break
		}
	}
	s.NotNil(found, "после ухода bob из команды задача должна всплыть")
	s.Equal(teamID, found.TeamID)
	s.Equal(bobUser.ID, found.AssigneeID)
	s.Equal(bobUser.Email, found.AssigneeEmail)
}
