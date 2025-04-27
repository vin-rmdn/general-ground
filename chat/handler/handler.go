package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vin-rmdn/general-ground/chat"
)

type service interface {
	Get(ctx context.Context, to string) ([]chat.Chat, error)
	Chat(ctx context.Context, to, message string) error
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{
		service: service,
	}
}

func (h handler) Get(c echo.Context) error {
	r := c.Request()

	ctx := r.Context()
	userID := r.Header.Get("User-ID")
	if userID == "" {
		return c.String(http.StatusBadRequest, "missing 'User-ID' header")
	}

	ctx = context.WithValue(ctx, chat.FromKey{}, userID)

	otherUserID := r.URL.Query().Get("with")
	if otherUserID == "" {
		return c.String(http.StatusBadRequest, "missing 'with' query parameter")
	}

	chats, err := h.service.Get(ctx, otherUserID)
	if err != nil {
		_ = c.String(http.StatusInternalServerError, fmt.Sprintf("failed to get chats: %v", err))

		return fmt.Errorf("failed to get chats: %w", err)
	}

	c.Response().Header().Add("Content-Type", "application/json")
	_ = c.JSON(http.StatusOK, chats)

	return nil
}

type chatRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (h handler) Chat(c echo.Context) error {
	r := c.Request()

	ctx := r.Context()
	userID := r.Header.Get("User-ID")
	if userID == "" {
		return c.String(http.StatusBadRequest, "missing 'User-ID' header")
	}

	ctx = context.WithValue(ctx, chat.FromKey{}, userID)

	var payload chatRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return c.String(http.StatusBadRequest, "failed to decode request body")
	}

	if err := h.service.Chat(ctx, payload.To, payload.Message); err != nil {
		return c.String(http.StatusInternalServerError, "failed to send chat")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	_ = c.String(http.StatusCreated, "message sent successfully")

	return nil
}
