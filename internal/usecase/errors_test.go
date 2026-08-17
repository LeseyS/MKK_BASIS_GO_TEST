package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errDB = errors.New("database is unreachable")

func TestRequireMembership_DBFailureIsNotForbidden(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: 1, Title: "task",
	})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrForbidden,
		"сбой БД не должен выглядеть как отказ в доступе")
}

func TestCreateTask_WriteFailure(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		createTask: func(context.Context, domain.Task) (domain.Task, error) {
			return domain.Task{}, errDB
		},
	}
	rd := &mockRedis{}
	uc := newUseCase(my, rd, nil, nil)

	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: 1, Title: "task",
	})

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, rd.invalidatedTeams, "несозданная задача не должна сбрасывать кеш")
}

func TestListTasks_MembershipCheckDBFailure(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrForbidden)
}

func TestListTasks_DBFailure(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			return nil, 0, errDB
		},
	}
	rd := &mockRedis{
		getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
			return dto.ListTasksOut{}, apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	_, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, rd.setFilters, "неудачную выборку в кеш не кладём")
}

func TestListTasks_CacheWriteFailureIsNotFatal(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			return []domain.Task{{ID: 1}}, 1, nil
		},
	}
	rd := &mockRedis{
		getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
			return dto.ListTasksOut{}, apperr.ErrNotFound
		},
		setTaskList: func(context.Context, dto.ListTasks, dto.ListTasksOut) error {
			return errors.New("redis is down")
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	out, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.NoError(t, err)
	assert.Equal(t, int64(1), out.Total)
}

func TestUpdateTask_ReadFailure(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{}, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{TaskID: 1, ActorID: 1})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrNotFound)
}

func TestUpdateTask_RereadFailure(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			return nil
		},
	}
	my.getTaskByID = func(context.Context, int64) (domain.Task, error) {
		if my.getTaskByIDCalls > 1 {
			return domain.Task{}, errDB
		}
		return domain.Task{ID: 1, TeamID: 7}, nil
	}
	rd := &mockRedis{}
	uc := newUseCase(my, rd, nil, nil)

	status := "done"
	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: 1, Status: &status,
	})

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, rd.invalidatedTeams)
}

func TestUpdateTask_TaskDisappearedDuringWrite(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			return apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	status := "done"
	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: 1, Status: &status,
	})

	require.ErrorIs(t, err, apperr.ErrNotFound)
}

func TestUpdateTask_CacheFailureDoesNotFailRequest(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			return nil
		},
	}
	rd := &mockRedis{
		invalidate: func(context.Context, int64) error {
			return errors.New("redis is down")
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	status := "done"
	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: 1, Status: &status,
	})

	require.NoError(t, err, "недоступный кеш не должен ронять уже выполненную правку")
}

func TestTaskHistory_ReadFailure(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{}, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{TaskID: 1, ActorID: 1})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrNotFound)
}

func TestTaskHistory_QueryFailure(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		taskHistory: func(context.Context, dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
			return nil, 0, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{TaskID: 1, ActorID: 1})

	require.ErrorIs(t, err, errDB)
}

func TestCreateWithOwner_TeamInsertFailure(t *testing.T) {
	my := &mockMySQL{
		createTeam: func(context.Context, dto.CreateTeamIn) (int64, error) {
			return 0, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateWithOwner(context.Background(), dto.CreateTeamIn{Name: "Alpha", OwnerID: 1})

	require.ErrorIs(t, err, errDB)
}

func TestTeamListForUser_DBFailure(t *testing.T) {
	my := &mockMySQL{
		teamListForUser: func(context.Context, int64) ([]domain.Team, error) {
			return nil, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.TeamListForUser(context.Background(), 1)

	require.ErrorIs(t, err, errDB)
}

func TestCreateUser_DBFailure(t *testing.T) {
	my := &mockMySQL{
		createUser: func(context.Context, domain.User) (int64, error) {
			return 0, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateUser(context.Background(), dto.CreateUserIn{
		Name: "alice", Email: "alice@example.com", Password: "password123",
	})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrEmailTaken)
}

func TestUserLogin_DBFailure(t *testing.T) {
	my := &mockMySQL{
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, errDB
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.UserLogin(context.Background(), dto.UserLoginIn{
		Email: "user@example.com", Password: "password123",
	})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrInvalidCredentials)
}

func TestInviteUser_TeamReadFailure(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.getTeamByID = func(context.Context, int64) (domain.Team, error) {
		return domain.Team{}, errDB
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, errDB)
}

func TestInviteUser_RoleLookupDBFailure(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.getMemberRole = func(context.Context, int64, int64) (domain.Role, error) {
		return "", errDB
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrForbidden)
}

func TestInviteUser_MemberInsertDBFailure(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.addMemberTeam = func(context.Context, dto.AddMemberTeam) error {
		return errDB
	}
	nt := &mockNotifier{}
	uc := newUseCase(my, nil, nil, nt)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, nt.sentTo, "письмо не должно уходить, если участник не добавлен")
}

func TestInviteUser_UserLookupDBFailure(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.getUserByEmail = func(context.Context, string) (domain.User, error) {
		return domain.User{}, errDB
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, apperr.ErrNotFound)
}
