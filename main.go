package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"
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
	r.HandleFunc("/", handler)
	r.HandleFunc("/calculate", calculateHandler)

	// start: set up any of your logger configuration here if necessary

	lf, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}

	log.Logger = zerolog.New(lf).With().Timestamp().Logger()
	log.Info().Msg("log init")

	// end: set up any of your logger configuration here

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
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

func handler(w http.ResponseWriter, r *http.Request) {
	// generate request id by uuid and add to context
	log := log.With().
		Str("request_id", uuid.New().String()).
		Logger()
	// creating a new context from the logger instance
	ctx := log.WithContext(r.Context())

	name := r.URL.Query().Get("name")
	res, err := greeting(ctx, name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func greeting(ctx context.Context, name string) (string, error) {
	if len(name) < 5 {
		log.Debug().
			// put function name here
			Str("function", "greeting").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Msg("name is too short: " + name)
		return fmt.Sprintf("Hello %s! Your name is to short\n", name), nil
	}
	return fmt.Sprintf("Hi %s", name), nil
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Str("request_id", uuid.New().String()).Logger()
	ctx := log.WithContext(r.Context())

	query := r.URL.Query()
	num1Str := query.Get("num1")
	num2Str := query.Get("num2")
	op := query.Get("op")

	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		log.Error().Str("function", "calculateHandler").Err(err).Msg("invalid num1 parameter")
		http.Error(w, "invalid num1 parameter", http.StatusBadRequest)
		return
	}

	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		log.Error().Str("function", "calculateHandler").Err(err).Msg("invalid num2 parameter")
		http.Error(w, "invalid num2 parameter", http.StatusBadRequest)
		return
	}

	result, err := calculate(ctx, num1, num2, op)
	if err != nil {
		log.Error().Str("function", "calculateHandler").Err(err).Msg("calculation error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("%.2f", result)))
}

func calculate(ctx context.Context, num1, num2 float64, op string) (float64, error) {
	log := zerolog.Ctx(ctx)

	switch op {
	case "add":
		log.Debug().
			Str("function", "calculate").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Msgf("performing addition: %.2f + %.2f", num1, num2)
		return num1 + num2, nil
	case "sub":
		log.Debug().
			Str("function", "calculate").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Msgf("performing subtraction: %.2f - %.2f", num1, num2)
		return num1 - num2, nil
	case "mult":
		log.Debug().
			Str("function", "calculate").
			Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
			Msgf("performing multiplication: %.2f * %.2f", num1, num2)
		return num1 * num2, nil
	case "div":
		if num2 == 0 {
			log.Err(fmt.Errorf("division by zero")).
				Str("function", "calculate").
				Str("request_id", fmt.Sprint(ctx.Value("request_id"))).
				Msg("division by zero")
			return 0, fmt.Errorf("division by zero")
		}
		log.Debug().Str("function", "calculate").Msgf("performing division: %.2f / %.2f", num1, num2)
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operation: %s", op)
	}
}
