package server

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"

	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/general-ground/chat/handler"
	"github.com/vin-rmdn/general-ground/chat/repository"
	"github.com/vin-rmdn/general-ground/chat/service"
)

func NewServer() {
	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong\n"))
	}))

	chatRepository := repository.New()
	chatService := service.New(chatRepository)
	chatHandler := handler.New(chatService)

	mux.HandleFunc("GET /chat", chatHandler.Get)
	mux.HandleFunc("POST /chat", chatHandler.Chat)

	cert, err := tls.LoadX509KeyPair(
		os.Getenv("CERTIFICATE_PATH"),
		os.Getenv("KEY_PATH"),
	)
	if err != nil {
		slog.Error("Failed to load certificate", "error", err)
		panic(err)
	}

	server := http3.Server{
		Addr: "0.0.0.0:8080",
		Port: 8080,
		TLSConfig: http3.ConfigureTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}),
		Handler:         mux,
		EnableDatagrams: true,
		Logger:          slog.Default(),
	}

	slog.Debug("Starting server on :8080")

	if err = server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)

		panic(err)
	}
}
