package handler

import (
	"context"
	"encoding/json"
	"net/http"

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

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "missing 'User-ID' header", http.StatusBadRequest)

		return
	}

	ctx = context.WithValue(ctx, chat.FromKey{}, userID)

	otherUserID := r.URL.Query().Get("with")
	if otherUserID == "" {
		http.Error(w, "missing 'with' query parameter", http.StatusBadRequest)

		return
	}

	chats, err := h.service.Get(ctx, otherUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	payload, err := json.Marshal(chats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to marshal response"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(payload)
}

type chatRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (h handler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "missing 'User-ID' header", http.StatusBadRequest)

		return
	}

	ctx = context.WithValue(ctx, chat.FromKey{}, userID)

	var payload chatRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)

		return
	}

	if err := h.service.Chat(ctx, payload.To, payload.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte("message sent successfully"))
}
