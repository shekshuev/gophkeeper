package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/handler"

	"github.com/shekshuev/gophkeeper/internal/repository"
	"github.com/shekshuev/gophkeeper/internal/service"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func NewServer(cfg *config.Config) *http.Server {
	userRepo := repository.NewUserRepositoryImpl(cfg)
	secretRepo := repository.NewSecretRepositoryImpl(cfg)
	userService := service.NewUserServiceImpl(userRepo, cfg)
	authService := service.NewAuthServiceImpl(userRepo, cfg)
	secretService := service.NewSecretServiceImpl(secretRepo)
	userHandler := handler.NewHandler(userService, authService, secretService, cfg)

	return &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: userHandler.Router,
	}
}

func main() {
	printBuildInfo()
	cfg := config.GetConfig()
	server := NewServer(&cfg)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server")
		}
	}()
	log.Print("Server listening on ", cfg.ServerAddress)
	<-done
	log.Print("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	} else {
		log.Print("Server shutdown gracefully")
	}
}
