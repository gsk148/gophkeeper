package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/config"
	"github.com/gsk148/gophkeeper/internal/app/server/handlers"
	"github.com/gsk148/gophkeeper/internal/app/server/storage"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	printCompilationInfo()
	cfg := config.MustLoad()
	s := getServer(cfg)
	idleConnectionsClosed := make(chan any)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
		<-exit
		stopServer(s)
		close(idleConnectionsClosed)
	}()

	go startServer(s)
	<-idleConnectionsClosed
}

func getServer(cfg *config.Config) *http.Server {
	db, err := storage.NewStorage(
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Name,
			cfg.Password,
			cfg.Sslmode))
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.NewHandler(db)
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: h,
	}
}

func startServer(s *http.Server) {
	if err := s.ListenAndServeTLS("internal/app/cert/server.crt", "internal/app/cert/server.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

func stopServer(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

func printCompilationInfo() {
	version := getCompilationInfoValue(buildVersion)
	date := getCompilationInfoValue(buildDate)
	fmt.Printf("Build version: %s\nBuild date: %s\n\n", version, date)
}

func getCompilationInfoValue(v string) string {
	if v != "" {
		return v
	}
	return "N/A"
}
