package usecase

import (
	"context"
	"fmt"

	"logging-challenge/domain"

	"github.com/rs/zerolog"
)

// CalculatorUseCase handles calculator business logic
type CalculatorUseCase struct {
	calculator *domain.Calculator
}

// NewCalculatorUseCase creates a new calculator use case instance
func NewCalculatorUseCase() *CalculatorUseCase {
	return &CalculatorUseCase{
		calculator: &domain.Calculator{},
	}
}

// Calculate performs arithmetic operations with logging
func (c *CalculatorUseCase) Calculate(ctx context.Context, num1, num2 float64, op string) (float64, error) {
	log := zerolog.Ctx(ctx)

	result, err := c.calculator.Calculate(num1, num2, op)
	if err != nil {
		log.Error().
			Str("function", "calculate").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Err(err).
			Msg("calculation error")
		return 0, err
	}

	log.Debug().
		Str("function", "calculate").
		Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
		Msgf("performing %s: %.2f %s %.2f = %.2f", op, num1, getOperatorSymbol(op), num2, result)

	return result, nil
}

// getOperatorSymbol returns the mathematical symbol for the operation
func getOperatorSymbol(op string) string {
	switch op {
	case "add":
		return "+"
	case "sub":
		return "-"
	case "mult":
		return "*"
	case "div":
		return "/"
	default:
		return op
	}
}
