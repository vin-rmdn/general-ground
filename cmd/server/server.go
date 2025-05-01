package server

import (
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/general-ground/chat/handler"
	"github.com/vin-rmdn/general-ground/chat/repository"
	"github.com/vin-rmdn/general-ground/chat/service"
	"github.com/vin-rmdn/general-ground/cmd/server/middleware"
)

type server struct {
	server http3.Server
}

func New(certificatePath, certificateKeyPath string) (*server, error) {
	router := echo.New()
	router.Use(middleware.Logger, middleware.Recovery)

	dependencies := setupServerDependencies()

	router.GET("/ping", ping)
	router.GET("/chat", dependencies.chatHandler.Get)
	router.POST("/chat", dependencies.chatHandler.Chat)

	cert, err := tls.LoadX509KeyPair(certificatePath, certificateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	return &server{
		server: http3.Server{
			Addr: "0.0.0.0:8080",
			Port: 8080,
			TLSConfig: http3.ConfigureTLSConfig(&tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS13,
			}),
			Handler:         router,
			EnableDatagrams: true,
			Logger:          dependencies.logger,
		},
	}, nil
}

func (s *server) Start() error {
	slog.Debug("Starting server on :8080")

	if err := s.server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)

		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

type chatHandler interface {
	Get(c echo.Context) error
	Chat(c echo.Context) error
}

type dependencies struct {
	chatHandler chatHandler
	logger      *slog.Logger
}

func setupServerDependencies() dependencies {
	chatRepository := repository.New()
	chatService := service.New(chatRepository)

	return dependencies{
		chatHandler: handler.New(chatService),
		logger:      slog.Default(),
	}
}
