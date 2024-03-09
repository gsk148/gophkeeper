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

	"github.com/gsk148/gophkeeper/internal/app/server/config"
	"github.com/gsk148/gophkeeper/internal/app/server/handlers"
)

var (
	buildVersion string
	buildDate    string
)

type ServerConfig interface {
	GetRepoURL() string
	GetServerAddress() string
	IsServerSecure() bool
}

func main() {
	printCompilationInfo()
	sCfg := config.MustLoad()
	s, err := getServer(sCfg)
	if err != nil {
		log.Fatal(err)
	}
	idleConnectionsClosed := make(chan any)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
		<-exit
		stopServer(s)
		close(idleConnectionsClosed)
	}()

	go startServer(s, sCfg)
	<-idleConnectionsClosed
}

func getServer(sCfg ServerConfig) (*http.Server, error) {
	h, err := handlers.NewHandler(sCfg.GetRepoURL())
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:              sCfg.GetServerAddress(),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s, nil
}

func startServer(s *http.Server, sCfg ServerConfig) {
	if sCfg.IsServerSecure() {
		err := s.ListenAndServeTLS("internal/app/cert/server.crt", "internal/app/cert/server.key")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	} else {
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
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
