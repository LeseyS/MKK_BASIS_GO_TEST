//go:build integration

package test

import (
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"
)

func (s *Suite) TestRegisterAndLogin() {
	email := fmt.Sprintf("register-%d@example.com", userCounter.Add(1))

	created, err := s.api.Register(ctx, task_client.RegisterRequest{
		Name: "alice", Email: email, Password: testPassword,
	})
	s.NoError(err)
	s.NotZero(created.ID)

	out, err := s.api.Login(ctx, task_client.LoginRequest{Email: email, Password: testPassword})
	s.NoError(err)
	s.NotEmpty(out.Token)
	s.Equal(email, out.User.Email)
	s.Equal(created.ID, out.User.ID)
}

func (s *Suite) TestRegisterDuplicateEmail() {
	email := fmt.Sprintf("duplicate-%d@example.com", userCounter.Add(1))

	_, err := s.api.Register(ctx, task_client.RegisterRequest{
		Name: "alice", Email: email, Password: testPassword,
	})
	s.NoError(err)

	_, err = s.api.Register(ctx, task_client.RegisterRequest{
		Name: "bob", Email: email, Password: testPassword,
	})
	s.ErrorIs(err, task_client.ErrConflict)
}

func (s *Suite) TestRegisterValidation() {
	cases := []struct {
		name string
		in   task_client.RegisterRequest
	}{
		{"короткое имя", task_client.RegisterRequest{Name: "ab", Email: "a@example.com", Password: testPassword}},
		{"кривой email", task_client.RegisterRequest{Name: "alice", Email: "not-an-email", Password: testPassword}},
		{"короткий пароль", task_client.RegisterRequest{Name: "alice", Email: "a@example.com", Password: "short"}},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			_, err := s.api.Register(ctx, c.in)
			s.ErrorIs(err, task_client.ErrBadRequest)
		})
	}
}

func (s *Suite) TestLoginWrongPassword() {
	_, user := s.newUser("wrongpass")

	_, err := s.api.Login(ctx, task_client.LoginRequest{
		Email: user.Email, Password: "not-the-password",
	})
	s.ErrorIs(err, task_client.ErrUnauthorized)
}

func (s *Suite) TestLoginUnknownEmail() {
	_, err := s.api.Login(ctx, task_client.LoginRequest{
		Email: "nobody@example.com", Password: testPassword,
	})
	s.ErrorIs(err, task_client.ErrUnauthorized)
}

func (s *Suite) TestRequestsWithoutTokenAreRejected() {
	_, err := s.api.Teams(ctx)
	s.ErrorIs(err, task_client.ErrUnauthorized)

	_, err = s.api.CreateTeam(ctx, task_client.CreateTeamRequest{Name: "Alpha"})
	s.ErrorIs(err, task_client.ErrUnauthorized)

	_, err = s.api.ListTasks(ctx, task_client.ListTasksRequest{TeamID: 1})
	s.ErrorIs(err, task_client.ErrUnauthorized)
}

func (s *Suite) TestGarbageTokenIsRejected() {
	broken := s.api.WithToken("not.a.jwt")

	_, err := broken.Teams(ctx)
	s.ErrorIs(err, task_client.ErrUnauthorized)
}
