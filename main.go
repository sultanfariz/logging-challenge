package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"logging-challenge/handler"
	"logging-challenge/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		oscall := <-ch
		log.Warn().Msgf("system call:%+v", oscall)
		cancel()
	}()

	r := mux.NewRouter()
	calculatorHandler := handler.NewCalculatorHandler()
	greetingHandler := handler.NewGreetingHandler()
	r.HandleFunc("/calculate", calculatorHandler.CalculateHandler)
	r.HandleFunc("/greet", greetingHandler.GreetHandler)

	// Wrap router with logger middleware
	handler := middleware.LoggerMiddleware(r)

	lf, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}

	log.Logger = zerolog.New(lf).With().Timestamp().Logger()
	log.Info().Msg("log init")

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen and serve http server")
		}
	}()
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}
}
