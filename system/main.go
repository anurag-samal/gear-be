package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github/anurag/altar-be/system/constants"
	"github/anurag/altar-be/system/logger"
	"github/anurag/altar-be/system/server"
)

func main() {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	logger.InitDefault(env)

	srv, err := server.NewServer()
	if err != nil {
		slog.Error("server initialization failed", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("server runtime error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	srv.ShutdownWithTimeout(constants.FORCE_SHUTDOWN_TIMEOUT)
}
