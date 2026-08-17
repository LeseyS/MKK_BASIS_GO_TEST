//go:build integration

package test

import "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"

func (s *Suite) TestRateLimitIsPerUser() {
	alice, _ := s.newUser("ratealice")
	bob, _ := s.newUser("ratebob")

	var lastErr error
	for i := 0; i < rateLimit+1; i++ {
		if _, lastErr = alice.Teams(ctx); lastErr != nil {
			break
		}
	}
	s.ErrorIs(lastErr, task_client.ErrTooManyRequests, "alice должна упереться в собственный лимит")

	_, err := bob.Teams(ctx)
	s.NoError(err, "исчерпанный лимит alice не должен задевать bob при общем IP")
}
