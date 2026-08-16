package email_service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/sony/gobreaker"
)

type EmailService struct {
	breaker          *gobreaker.CircuitBreaker
	failRate         float64
	simulatedLatency time.Duration
}

func NewEmailService() *EmailService {
	settings := gobreaker.Settings{
		Name:        "email-service",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// trip after 5 consecutive failures
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Warn("circuit breaker state change", "breaker", name, "from", from.String(), "to", to.String())
		},
	}
	return &EmailService{
		breaker:          gobreaker.NewCircuitBreaker(settings),
		failRate:         0.15,
		simulatedLatency: 50 * time.Millisecond,
	}
}

func (s *EmailService) SendInvite(ctx context.Context, toEmail, teamName string) error {
	_, err := s.breaker.Execute(func() (interface{}, error) {
		return nil, s.mockSend(ctx, toEmail, teamName)
	})
	return err
}

func (s *EmailService) mockSend(ctx context.Context, toEmail, teamName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.simulatedLatency):
	}

	if rand.Float64() < s.failRate {
		return fmt.Errorf("mock email provider: temporary failure sending invite to %s", toEmail)
	}

	slog.Info("mock invite email sent", "to", toEmail, "team", teamName)
	return nil
}
