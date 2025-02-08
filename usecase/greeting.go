package usecase

import (
	"context"
	"fmt"

	"logging-challenge/domain"

	"github.com/rs/zerolog"
)

// GreetingUseCase handles greeting business logic
type GreetingUseCase struct {
	greeting *domain.Greeting
}

// NewGreetingUseCase creates a new greeting use case instance
func NewGreetingUseCase() *GreetingUseCase {
	return &GreetingUseCase{
		greeting: &domain.Greeting{},
	}
}

// GetGreeting generates a greeting message with logging
func (g *GreetingUseCase) GetGreeting(ctx context.Context, name string) (string, error) {
	log := zerolog.Ctx(ctx)

	result, err := g.greeting.GetGreeting(name)
	if err != nil {
		log.Error().
			Str("function", "getGreeting").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Err(err).
			Msg("greeting error")
		return "", err
	}

	log.Debug().
		Str("function", "getGreeting").
		Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
		Str("name", name).
		Str("result", result).
		Msg("greeting generated")

	return result, nil
}
