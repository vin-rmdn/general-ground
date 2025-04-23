package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vin-rmdn/general-ground/chat"
	"github.com/vin-rmdn/general-ground/chat/handler"
)

type mockService struct {
	GetFn  func(ctx context.Context, to string) ([]chat.Chat, error)
	ChatFn func(ctx context.Context, to, message string) error
}

func (m *mockService) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	return m.GetFn(ctx, to)
}

func (m *mockService) Chat(ctx context.Context, to, message string) error {
	return m.ChatFn(ctx, to, message)
}

func TestHandler_Get(t *testing.T) {
	t.Run("should return 400 when 'User-ID' header is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?with=user2", nil)
		rec := httptest.NewRecorder()

		h := handler.New(&mockService{})
		h.Get(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}

		expectedBody := "missing 'User-ID' header\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 400 when 'with' query parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		h := handler.New(&mockService{})
		h.Get(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}

		expectedBody := "missing 'with' query parameter\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 500 when service.Get returns an error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?with=user2", nil)
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		mockSvc := &mockService{
			GetFn: func(ctx context.Context, to string) ([]chat.Chat, error) {
				return nil, errors.New("service error")
			},
		}

		h := handler.New(mockSvc)
		h.Get(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", rec.Code)
		}

		expectedBody := "service error\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 200 with chats when service.Get succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?with=user2", nil)
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		mockChats := []chat.Chat{
			{From: "user1", To: "user2", Message: "Hello, user2!"},
		}

		mockSvc := &mockService{
			GetFn: func(ctx context.Context, to string) ([]chat.Chat, error) {
				return mockChats, nil
			},
		}

		h := handler.New(mockSvc)
		h.Get(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		expectedBody := `[{"from":"user1","to":"user2","message":"Hello, user2!","timestamp":"0001-01-01T00:00:00Z"}]`
		if strings.TrimSpace(rec.Body.String()) != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}

		if rec.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected Content-Type 'application/json', got '%s'", rec.Header().Get("Content-Type"))
		}
	})
}

func TestHandler_Chat(t *testing.T) {
	t.Run("should return 400 when 'User-ID' header is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"to":"user2","message":"Hello!"}`))
		rec := httptest.NewRecorder()

		h := handler.New(&mockService{})
		h.Chat(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}

		expectedBody := "missing 'User-ID' header\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid-json`))
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		h := handler.New(&mockService{})
		h.Chat(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}

		expectedBody := "failed to decode request body\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 500 when service.Chat returns an error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"to":"user2","message":"Hello!"}`))
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		mockSvc := &mockService{
			ChatFn: func(ctx context.Context, to, message string) error {
				return errors.New("service error")
			},
		}

		h := handler.New(mockSvc)
		h.Chat(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", rec.Code)
		}

		expectedBody := "service error\n"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("should return 201 when service.Chat succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"to":"user2","message":"Hello!"}`))
		req.Header.Set("User-ID", "user1")
		rec := httptest.NewRecorder()

		mockSvc := &mockService{
			ChatFn: func(ctx context.Context, to, message string) error {
				return nil
			},
		}

		h := handler.New(mockSvc)
		h.Chat(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rec.Code)
		}

		expectedBody := "message sent successfully"
		if rec.Body.String() != expectedBody {
			t.Fatalf("expected body '%s', got '%s'", expectedBody, rec.Body.String())
		}

		if rec.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected Content-Type 'application/json', got '%s'", rec.Header().Get("Content-Type"))
		}
	})
}
