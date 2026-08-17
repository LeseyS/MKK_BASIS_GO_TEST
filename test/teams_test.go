//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestCreateTeamMakesCreatorOwner() {
	alice, aliceUser := s.newUser("owner")

	out, err := alice.CreateTeam(ctx, task_client.CreateTeamRequest{Name: "Alpha"})
	s.NoError(err)
	s.NotZero(out.Team.ID)
	s.Equal("Alpha", out.Team.Name)
	s.Equal(aliceUser.ID, out.Team.CreatedBy)

	teams, err := alice.Teams(ctx)
	s.NoError(err)
	s.Len(teams, 1)
	s.Equal(out.Team.ID, teams[0].Team.ID)
}

func (s *Suite) TestCreateTeamValidation() {
	alice, _ := s.newUser("teamvalidation")

	_, err := alice.CreateTeam(ctx, task_client.CreateTeamRequest{Name: ""})
	s.ErrorIs(err, task_client.ErrBadRequest)
}

func (s *Suite) TestTeamListShowsOnlyOwnTeams() {
	alice, _ := s.newUser("listowner")
	bob, _ := s.newUser("listoutsider")

	s.newTeam(alice, "Alpha")

	teams, err := bob.Teams(ctx)
	s.NoError(err)
	s.Empty(teams, "чужие команды в списке появляться не должны")
}

func (s *Suite) TestInviteAddsMember() {
	alice, _ := s.newUser("inviter")
	bob, bobUser := s.newUser("invitee")

	teamID := s.newTeam(alice, "Alpha")

	out, err := alice.InviteUser(ctx, teamID, task_client.InviteRequest{Email: bobUser.Email})
	s.NoError(err)
	s.True(out.Invited)

	teams, err := bob.Teams(ctx)
	s.NoError(err)
	s.Len(teams, 1)
	s.Equal(teamID, teams[0].Team.ID)
}

func (s *Suite) TestInviteTwiceIsConflict() {
	alice, _ := s.newUser("inviter2")
	_, bobUser := s.newUser("invitee2")

	teamID := s.newTeam(alice, "Alpha")
	s.invite(alice, teamID, bobUser.Email)

	_, err := alice.InviteUser(ctx, teamID, task_client.InviteRequest{Email: bobUser.Email})
	s.ErrorIs(err, task_client.ErrConflict)
}

func (s *Suite) TestInviteUnknownEmail() {
	alice, _ := s.newUser("inviter3")
	teamID := s.newTeam(alice, "Alpha")

	_, err := alice.InviteUser(ctx, teamID, task_client.InviteRequest{
		Email: "definitely-not-registered@example.com",
	})
	s.ErrorIs(err, task_client.ErrNotFound)
}

func (s *Suite) TestOrdinaryMemberCannotInvite() {
	alice, _ := s.newUser("inviter4")
	bob, bobUser := s.newUser("member4")
	_, carolUser := s.newUser("outsider4")

	teamID := s.newTeam(alice, "Alpha")
	s.invite(alice, teamID, bobUser.Email)

	_, err := bob.InviteUser(ctx, teamID, task_client.InviteRequest{Email: carolUser.Email})
	s.ErrorIs(err, task_client.ErrForbidden, "приглашать может только owner или admin")
}

func (s *Suite) TestOutsiderCannotInvite() {
	alice, _ := s.newUser("inviter5")
	bob, _ := s.newUser("outsider5")
	_, carolUser := s.newUser("target5")

	teamID := s.newTeam(alice, "Alpha")

	_, err := bob.InviteUser(ctx, teamID, task_client.InviteRequest{Email: carolUser.Email})
	s.ErrorIs(err, task_client.ErrForbidden)
}
