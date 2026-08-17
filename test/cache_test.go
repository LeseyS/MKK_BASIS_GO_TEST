//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestTaskListIsCachedAndInvalidated() {
	alice, aliceUser := s.newUser("cacheowner")
	teamID := s.newTeam(alice, "Cache")
	s.newTask(alice, teamID, "первая")

	first, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(1), first.Total)

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO tasks (team_id, title, description, status, created_by, created_at, updated_at)
		 VALUES (?, ?, '', 'todo', ?, NOW(), NOW())`, teamID, "вставлена мимо API", aliceUser.ID)
	s.NoError(err)

	cached, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(1), cached.Total, "ответ должен прийти из кеша и не увидеть прямую вставку")

	s.newTask(alice, teamID, "третья")

	fresh, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(3), fresh.Total, "создание задачи сбрасывает кеш команды")
}

func (s *Suite) TestUpdateInvalidatesCache() {
	alice, _ := s.newUser("cacheupdate")
	teamID := s.newTeam(alice, "CacheUpdate")
	task := s.newTask(alice, teamID, "задача")

	before, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal("todo", before.Tasks[0].Status)

	done := "done"
	_, err = alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &done})
	s.NoError(err)

	after, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal("done", after.Tasks[0].Status, "правка задачи должна сбрасывать кеш списка")
}

func (s *Suite) TestCacheIsPerFilterCombination() {
	alice, _ := s.newUser("cachefilters")
	teamID := s.newTeam(alice, "CacheFilters")
	task := s.newTask(alice, teamID, "задача")

	all, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID})
	s.NoError(err)
	s.Equal(int64(1), all.Total)

	doneOnly, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID, Status: "done"})
	s.NoError(err)
	s.Equal(int64(0), doneOnly.Total, "фильтр по статусу не должен переиспользовать чужой ключ кеша")

	done := "done"
	_, err = alice.UpdateTask(ctx, task.ID, task_client.UpdateTaskRequest{Status: &done})
	s.NoError(err)

	doneAfter, err := alice.ListTasks(ctx, task_client.ListTasksRequest{TeamID: teamID, Status: "done"})
	s.NoError(err)
	s.Equal(int64(1), doneAfter.Total, "инвалидация должна затрагивать все ключи команды")
}
