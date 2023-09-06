package server

import (
	"context"
	"fmt"
	"gcipher/internal/certificate"
	"gcipher/internal/config"
	"gcipher/internal/db/repositories"
	ocsp "gcipher/internal/oscp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer initializes and starts the HTTP server.
func StartServer() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		return
	}

	// Initialize repositories
	if err := repositories.InitializeRepositories(); err != nil {
		log.Fatal("Failed to initialize repositories:", err)
	}

	// Create a new server mux
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/certificate/request", certificate.HandleCertificateRequest)
	mux.HandleFunc("/api/v1/certificate/retrieve", certificate.HandleCertificateRetrieval)
	mux.HandleFunc("/api/v1/certificate/revoke", certificate.HandleRevokeCertificate)
	mux.HandleFunc("/api/v1/certificate/list", certificate.HandleCertificateList)
	mux.HandleFunc("/public/ca/intermediate/crl", ocsp.HandleCRL)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}
	fmt.Println("Server gracefully stopped")
}
