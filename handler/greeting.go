package handler

import (
	"fmt"
	"net/http"

	"logging-challenge/middleware"
	"logging-challenge/usecase"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// GreetingHandler handles HTTP requests for greeting operations
type GreetingHandler struct {
	greetingUseCase *usecase.GreetingUseCase
}

// NewGreetingHandler creates a new greeting handler instance
func NewGreetingHandler() *GreetingHandler {
	return &GreetingHandler{
		greetingUseCase: usecase.NewGreetingUseCase(),
	}
}

// GreetHandler handles greeting requests
func (h *GreetingHandler) GreetHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Str("request_id", uuid.New().String()).Logger()
	ctx := log.WithContext(r.Context())

	name := r.URL.Query().Get("name")

	result, err := h.greetingUseCase.GetGreeting(ctx, name)
	if err != nil {
		middleware.LogError(ctx, "greetHandler", err, "greeting error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("%s", result)))
}
