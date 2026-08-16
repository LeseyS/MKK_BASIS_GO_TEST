package usecase

import (
	"context"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

//go:generate mockery

type Redis interface {
	InvalidateTeamTasks(ctx context.Context, teamID int64) error
	SetTaskList(ctx context.Context, f dto.ListTasks, out dto.ListTasksOut) error
	GetTaskList(ctx context.Context, f dto.ListTasks) (dto.ListTasksOut, error)
}

type MySQL interface {
	CreateUser(ctx context.Context, user domain.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)

	CreateTeam(ctx context.Context, in dto.CreateTeamIn) (int64, error)
	CreateTeamMembers(ctx context.Context, teamID, userID int64, role domain.Role) error
	TeamListForUser(ctx context.Context, userID int64) ([]domain.Team, error)
	GetMemberRole(ctx context.Context, teamID, userID int64) (domain.Role, error)
	AddMemberTeam(ctx context.Context, in dto.AddMemberTeam) error
	GetTeamByID(ctx context.Context, teamID int64) (domain.Team, error)

	CreateTask(ctx context.Context, t domain.Task) (domain.Task, error)
	ListTasks(ctx context.Context, f dto.ListTasks) ([]domain.Task, int64, error)
	GetTaskByID(ctx context.Context, id int64) (domain.Task, error)
	UpdateTask(ctx context.Context, in dto.UpdateTaskIn) error
}

type TokenIssuer interface {
	Generate(userID int64, email string) (string, error)
}

type Notifier interface {
	SendInvite(ctx context.Context, toEmail, teamName string) error
}

type UseCase struct {
	mysql MySQL
	redis Redis

	tokenIssuer TokenIssuer
	notifier    Notifier
}

func New(mysql MySQL, redis Redis, tokenIssuer TokenIssuer, notifier Notifier) *UseCase {
	return &UseCase{
		mysql:       mysql,
		redis:       redis,
		tokenIssuer: tokenIssuer,
		notifier:    notifier,
	}
}
