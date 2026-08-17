package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func TestMain(m *testing.M) {
	transaction.IsUnitTest = true

	os.Exit(m.Run())
}

type mockMySQL struct {
	MySQL

	createUser        func(context.Context, domain.User) (int64, error)
	getUserByEmail    func(context.Context, string) (domain.User, error)
	createTeam        func(context.Context, dto.CreateTeamIn) (int64, error)
	createTeamMembers func(context.Context, int64, int64, domain.Role) error
	teamListForUser   func(context.Context, int64) ([]domain.Team, error)
	getMemberRole     func(context.Context, int64, int64) (domain.Role, error)
	addMemberTeam     func(context.Context, dto.AddMemberTeam) error
	getTeamByID       func(context.Context, int64) (domain.Team, error)
	createTask        func(context.Context, domain.Task) (domain.Task, error)
	listTasks         func(context.Context, dto.ListTasks) ([]domain.Task, int64, error)
	getTaskByID       func(context.Context, int64) (domain.Task, error)
	updateTask        func(context.Context, dto.UpdateTaskIn) error
	taskHistory       func(context.Context, dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error)

	getTaskByIDCalls int
}

func (m *mockMySQL) CreateUser(ctx context.Context, u domain.User) (int64, error) {
	return m.createUser(ctx, u)
}

func (m *mockMySQL) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.getUserByEmail(ctx, email)
}

func (m *mockMySQL) CreateTeam(ctx context.Context, in dto.CreateTeamIn) (int64, error) {
	return m.createTeam(ctx, in)
}

func (m *mockMySQL) CreateTeamMembers(ctx context.Context, teamID, userID int64, role domain.Role) error {
	return m.createTeamMembers(ctx, teamID, userID, role)
}

func (m *mockMySQL) TeamListForUser(ctx context.Context, userID int64) ([]domain.Team, error) {
	return m.teamListForUser(ctx, userID)
}

func (m *mockMySQL) GetMemberRole(ctx context.Context, teamID, userID int64) (domain.Role, error) {
	return m.getMemberRole(ctx, teamID, userID)
}

func (m *mockMySQL) AddMemberTeam(ctx context.Context, in dto.AddMemberTeam) error {
	return m.addMemberTeam(ctx, in)
}

func (m *mockMySQL) GetTeamByID(ctx context.Context, teamID int64) (domain.Team, error) {
	return m.getTeamByID(ctx, teamID)
}

func (m *mockMySQL) CreateTask(ctx context.Context, t domain.Task) (domain.Task, error) {
	return m.createTask(ctx, t)
}

func (m *mockMySQL) ListTasks(ctx context.Context, f dto.ListTasks) ([]domain.Task, int64, error) {
	return m.listTasks(ctx, f)
}

func (m *mockMySQL) GetTaskByID(ctx context.Context, id int64) (domain.Task, error) {
	m.getTaskByIDCalls++
	return m.getTaskByID(ctx, id)
}

func (m *mockMySQL) UpdateTask(ctx context.Context, in dto.UpdateTaskIn) error {
	return m.updateTask(ctx, in)
}

func (m *mockMySQL) TaskHistory(ctx context.Context, in dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
	return m.taskHistory(ctx, in)
}

type mockRedis struct {
	getTaskList func(context.Context, dto.ListTasks) (dto.ListTasksOut, error)
	setTaskList func(context.Context, dto.ListTasks, dto.ListTasksOut) error
	invalidate  func(context.Context, int64) error

	invalidatedTeams []int64
	setFilters       []dto.ListTasks
}

func (m *mockRedis) InvalidateTeamTasks(ctx context.Context, teamID int64) error {
	m.invalidatedTeams = append(m.invalidatedTeams, teamID)
	if m.invalidate == nil {
		return nil
	}
	return m.invalidate(ctx, teamID)
}

func (m *mockRedis) SetTaskList(ctx context.Context, f dto.ListTasks, out dto.ListTasksOut) error {
	m.setFilters = append(m.setFilters, f)
	if m.setTaskList == nil {
		return nil
	}
	return m.setTaskList(ctx, f, out)
}

func (m *mockRedis) GetTaskList(ctx context.Context, f dto.ListTasks) (dto.ListTasksOut, error) {
	return m.getTaskList(ctx, f)
}

type mockTokenIssuer struct {
	generate func(int64, string) (string, error)
}

func (m *mockTokenIssuer) Generate(userID int64, email string) (string, error) {
	return m.generate(userID, email)
}

type mockNotifier struct {
	sendInvite func(context.Context, string, string) error

	sentTo []string
}

func (m *mockNotifier) SendInvite(ctx context.Context, toEmail, teamName string) error {
	m.sentTo = append(m.sentTo, toEmail)
	if m.sendInvite == nil {
		return nil
	}
	return m.sendInvite(ctx, toEmail, teamName)
}

func newUseCase(my *mockMySQL, rd *mockRedis, ti *mockTokenIssuer, nt *mockNotifier) *UseCase {
	if my == nil {
		my = &mockMySQL{}
	}
	if rd == nil {
		rd = &mockRedis{}
	}
	if ti == nil {
		ti = &mockTokenIssuer{}
	}
	if nt == nil {
		nt = &mockNotifier{}
	}

	return New(my, rd, ti, nt)
}
