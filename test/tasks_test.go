//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestCreateTaskDefaults() {
	alice, aliceUser := s.newUser("taskowner")
	teamID := s.newTeam(alice, "Alpha")

	out, err := alice.CreateTask(ctx, task_client.CreateTaskRequest{
		TeamID: teamID, Title: "первая задача",
	})
	s.NoError(err)
	s.NotZero(out.Task.ID)
	s.Equal("todo", out.Task.Status, "статус по умолчанию — todo")
	s.Equal(aliceUser.ID, out.Task.CreatedBy)
	s.Empty(out.Task.Description)
	s.Nil(out.Task.AssigneeID)
}

func (s *Suite) TestCreateTaskValidation() {
	alice, _ := s.newUser("taskvalidation")
	teamID := s.newTeam(alice, "Alpha")

	_, err := alice.CreateTask(ctx, task_client.CreateTaskRequest{TeamID: teamID, Title: ""})
	s.ErrorIs(err, task_client.ErrBadRequest)

	_, err = alice.CreateTask(ctx, task_client.CreateTaskRequest{
		TeamID: teamID, Title: "задача", Status: "bogus",
	})
	s.ErrorIs(err, task_client.ErrBadRequest)
}

func (s *Suite) TestAssigneeMustBeTeamMember() {
	alice, _ := s.newUser("assignowner")
	_, outsider := s.newUser("assignoutsider")
	teamID := s.newTeam(alice, "Alpha")

	_, err := alice.CreateTask(ctx, task_client.CreateTaskRequest{
		TeamID: teamID, Title: "задача", AssigneeID: &outsider.ID,
	})
	s.ErrorIs(err, task_client.ErrUnprocessable, "исполнитель обязан состоять в команде")
}

func (s *Suite) TestListTasksFilters() {
	alice, aliceUser := s.newUser("filterowner")
	bob, bobUser := s.newUser("filtermember")
	teamID := s.newTeam(alice, "Alpha")
	s.invite(alice, teamID, bobUser.Email)

	first := s.newTask(alice, teamID, "todo-задача")
	second := s.newTask(alice, teamID, "done-задача")
	_ = s.newTask(bob, teamID, "задача от bob")

	done := "done"
	_, err := alice.UpdateTask(ctx, second.ID, task_client.UpdateTaskRequest{Status: &done})
	s.NoError(err)

	_, err = alice.UpdateTask(ctx, first.ID, task_client.UpdateTaskRequest{AssigneeID: &aliceUser.ID})
	s.NoError(err)

	all, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(3), all.Total)

	byStatus, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID, Status: "done"})
	s.NoError(err)
	s.Equal(int64(1), byStatus.Total)
	s.Equal(second.ID, byStatus.Tasks[0].ID)

	byAssignee, err := alice.ListTasks(ctx, task_client.ListTasksRequest{
		TeamID: teamID, AssigneeID: aliceUser.ID,
	})
	s.NoError(err)
	s.Equal(int64(1), byAssignee.Total)
	s.Equal(first.ID, byAssignee.Tasks[0].ID)

	empty, err := alice.ListTasks(ctx, task_client.ListTasksRequest{
		TeamID: teamID, Status: "in_progress",
	})
	s.NoError(err)
	s.Equal(int64(0), empty.Total)
	s.Empty(empty.Tasks)
}

func (s *Suite) TestListTasksPagination() {
	alice, _ := s.newUser("paging")
	teamID := s.newTeam(alice, "Alpha")

	for i := 0; i < 5; i++ {
		s.newTask(alice, teamID, "задача")
	}

	page, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID, Limit: 2})
	s.NoError(err)
	s.Equal(int64(5), page.Total, "total считает все записи, а не страницу")
	s.Len(page.Tasks, 2)
	s.Equal(2, page.Limit)

	second, err := alice.ListTasks(ctx, task_client.ListTasksRequest{
		TeamID: teamID, Limit: 2, Offset: 2,
	})
	s.NoError(err)
	s.Len(second.Tasks, 2)
	s.NotEqual(page.Tasks[0].ID, second.Tasks[0].ID, "страницы не должны пересекаться")

	last, err := alice.ListTasks(ctx, task_client.ListTasksRequest{
		TeamID: teamID, Limit: 2, Offset: 4,
	})
	s.NoError(err)
	s.Len(last.Tasks, 1)
}

func (s *Suite) TestListTasksLimitIsClamped() {
	alice, _ := s.newUser("clamp")
	teamID := s.newTeam(alice, "Alpha")
	s.newTask(alice, teamID, "задача")

	out, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID, Limit: 500})
	s.NoError(err)
	s.Equal(100, out.Limit, "запрошенный лимит выше максимума ужимается до 100")

	def, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(20, def.Limit, "без параметра действует лимит по умолчанию")
}

func (s *Suite) TestUpdateTaskFields() {
	alice, _ := s.newUser("updater")
	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "исходный заголовок")

	title := "новый заголовок"
	description := "новое описание"
	status := "in_progress"

	out, err := alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{
		Title: &title, Description: &description, Status: &status,
	})
	s.NoError(err)
	s.Equal(title, out.Task.Title)
	s.Equal(description, out.Task.Description)
	s.Equal(status, out.Task.Status)
	s.Equal(task.ID, out.Task.ID)
}

func (s *Suite) TestUpdateTaskValidation() {
	alice, _ := s.newUser("updatevalidation")
	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "задача")

	bogus := "bogus"
	_, err := alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &bogus})
	s.ErrorIs(err, task_client.ErrBadRequest)

	empty := ""
	_, err = alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Title: &empty})
	s.ErrorIs(err, task_client.ErrBadRequest)
}

func (s *Suite) TestUpdateMissingTask() {
	alice, _ := s.newUser("updatemissing")

	status := "done"
	_, err := alice.UpdateTask(ctx, 999999, task_client.UpdateTaskRequest{Status: &status})
	s.ErrorIs(err, task_client.ErrNotFound)
}

func (s *Suite) TestTaskHistoryRecordsEveryChange() {
	alice, aliceUser := s.newUser("historyowner")
	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "исходный заголовок")

	title := "переименована"
	status := "in_progress"
	_, err := alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{
		Title: &title, Status: &status,
	})
	s.NoError(err)

	done := "done"
	_, err = alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &done})
	s.NoError(err)

	history, err := alice.TaskHistory(ctx, task.ID, 0, 0)
	s.NoError(err)
	s.Equal(int64(3), history.Total, "две правки, одна из которых затронула два поля")

	latest := history.History[0]
	s.Equal("status", latest.Field)
	s.Equal("in_progress", *latest.OldValue)
	s.Equal("done", *latest.NewValue)
	s.Equal(aliceUser.ID, latest.ChangedBy)
	s.Equal(aliceUser.Email, latest.ChangedByEmail)

	fields := map[string]bool{}
	for _, e := range history.History {
		fields[e.Field] = true
	}
	s.True(fields["title"])
	s.True(fields["status"])
}

func (s *Suite) TestTaskHistoryIsEmptyForUntouchedTask() {
	alice, _ := s.newUser("historyempty")
	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "нетронутая")

	history, err := alice.TaskHistory(ctx, task.ID, 0, 0)
	s.NoError(err)
	s.Equal(int64(0), history.Total)
	s.Empty(history.History)
}

func (s *Suite) TestTaskHistoryPagination() {
	alice, _ := s.newUser("historypaging")
	teamID := s.newTeam(alice, "Alpha")
	task := s.newTask(alice, teamID, "задача")

	for _, status := range []string{"in_progress", "done", "todo"} {
		next := status
		_, err := alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &next})
		s.NoError(err)
	}

	page, err := alice.TaskHistory(ctx, task.ID, 2, 0)
	s.NoError(err)
	s.Equal(int64(3), page.Total)
	s.Len(page.History, 2)
	s.Equal(2, page.Limit)

	rest, err := alice.TaskHistory(ctx, task.ID, 2, 2)
	s.NoError(err)
	s.Len(rest.History, 1)
}

func (s *Suite) TestTaskHistoryMissingTask() {
	alice, _ := s.newUser("historymissing")

	_, err := alice.TaskHistory(ctx, 999999, 0, 0)
	s.ErrorIs(err, task_client.ErrNotFound)
}
