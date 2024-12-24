package main

import (
	pb "auth-service/gen/go/auth"
	"auth-service/internal/app"
	"auth-service/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.LoadConfig()
	log := setupLogger(cfg.Env)
	log.Info("starting auth-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	application := app.New(log, cfg.GRPC.Port, cfg.DatabaseUrl, cfg.RabbitMQ.Url, cfg.RabbitMQ.QueueName, cfg.TokenTTL)

	defer func(application *app.App) {
		err := application.Close()
		if err != nil {
			panic(err)
		}
	}(application)

	go application.GRPCSrv.MustRun()

	go func() {
		log.Info("starting REST gateway")
		if err := runRESTGateway(cfg, log); err != nil {
			log.Error("failed to run REST gateway", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", os.Signal.String(sign)))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func runRESTGateway(cfg *config.Config, log *slog.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()} // todo: исправить deprecated

	err := pb.RegisterAuthHandlerFromEndpoint(ctx, mux, "localhost:"+strconv.Itoa(cfg.GRPC.Port), opts)
	if err != nil {
		return err
	}

	// Обработчик для установки и удаления cookie
	cookieHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" && r.Method == "POST" {
			// Логика установки cookie при входе
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, r)

			var loginResp pb.LoginResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &loginResp); err == nil && loginResp.Token != "" {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    loginResp.Token,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
				})
			}

			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			_, err := w.Write(rec.Body.Bytes())
			if err != nil {
				log.Error("Error writing response body")
			}
		} else if r.URL.Path == "/api/v1/auth/logout" && r.Method == "DELETE" {
			// Логика удаления cookie при выходе
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Unix(0, 0), // Устанавливаем срок действия в прошлое
				MaxAge:   -1,
			})
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "Logged out successfully"}`))
			if err != nil {
				log.Error("Error writing logout response")
			}
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	allowedOrigins := "http://localhost:3000"
	if cfg.Env == "prod" {
		allowedOrigins = "http://localhost:80"
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(cookieHandler)

	log.Info("REST gateway is running on port: " + strconv.Itoa(cfg.Rest.Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Rest.Port), corsHandler)
}
