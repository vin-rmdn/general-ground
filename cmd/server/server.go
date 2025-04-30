package server

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/general-ground/chat/handler"
	"github.com/vin-rmdn/general-ground/chat/repository"
	"github.com/vin-rmdn/general-ground/chat/service"
	"github.com/vin-rmdn/general-ground/cmd/server/middleware"
)

func NewServer() {
	router := echo.New()

	router.Use(middleware.Logger, middleware.Recovery)

	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	chatRepository := repository.New()
	chatService := service.New(chatRepository)
	chatHandler := handler.New(chatService)

	router.GET("/chat", chatHandler.Get)
	router.POST("/chat", chatHandler.Chat)

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
		Handler:         router,
		EnableDatagrams: true,
		Logger:          slog.Default(),
	}

	slog.Debug("Starting server on :8080")

	if err = server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)

		panic(err)
	}
}
