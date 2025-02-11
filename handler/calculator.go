package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"logging-challenge/middleware"
	"logging-challenge/usecase"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CalculatorHandler handles HTTP requests for calculator operations
type CalculatorHandler struct {
	calculatorUseCase *usecase.CalculatorUseCase
}

// NewCalculatorHandler creates a new calculator handler instance
func NewCalculatorHandler() *CalculatorHandler {
	return &CalculatorHandler{
		calculatorUseCase: usecase.NewCalculatorUseCase(),
	}
}

// CalculateHandler handles calculation requests
func (h *CalculatorHandler) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Str("request_id", uuid.New().String()).Logger()
	ctx := log.WithContext(r.Context())

	query := r.URL.Query()
	num1Str := query.Get("num1")
	num2Str := query.Get("num2")
	op := query.Get("op")

	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		middleware.LogError(ctx, "calculateHandler", err, "invalid num1 parameter")
		http.Error(w, "invalid num1 parameter", http.StatusBadRequest)
		return
	}

	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		middleware.LogError(ctx, "calculateHandler", err, "invalid num2 parameter")
		http.Error(w, "invalid num2 parameter", http.StatusBadRequest)
		return
	}

	result, err := h.calculatorUseCase.Calculate(ctx, num1, num2, op)
	if err != nil {
		middleware.LogError(ctx, "calculateHandler", err, "calculation error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("%.2f", result)))
}
