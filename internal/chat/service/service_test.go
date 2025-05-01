package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vin-rmdn/general-ground/internal/chat"
	"github.com/vin-rmdn/general-ground/internal/chat/service"
)

type mockRepository struct {
	chats map[string][]chat.Chat
	err   error
}

func (m *mockRepository) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.chats[to], nil
}

func (m *mockRepository) Save(ctx context.Context, to, message string, timestamp time.Time) error {
	if m.err != nil {
		return m.err
	}
	m.chats[to] = append(m.chats[to], chat.Chat{Message: message, Timestamp: timestamp})
	return nil
}

func TestService_Get(t *testing.T) {
	t.Run("should return chats successfully", func(t *testing.T) {
		mockRepo := &mockRepository{
			chats: map[string][]chat.Chat{
				"user2": {
					{Message: "Hello, user2!", Timestamp: time.Now()},
				},
			},
		}
		svc := service.New(mockRepo)

		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		chats, err := svc.Get(ctx, "user2")
		if err != nil {
			t.Fatalf("failed to get chats: %v", err)
		}

		if len(chats) != 1 {
			t.Fatalf("expected 1 chat, got %d", len(chats))
		}

		if chats[0].Message != "Hello, user2!" {
			t.Fatalf("expected message 'Hello, user2!', got '%s'", chats[0].Message)
		}
	})

	t.Run("should return error when repository Get fails", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}
		svc := service.New(mockRepo)

		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		_, err := svc.Get(ctx, "user2")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		expectedErr := "failed to get chats: repository error"
		if err.Error() != expectedErr {
			t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("should return empty chats when no chats exist", func(t *testing.T) {
		mockRepo := &mockRepository{
			chats: map[string][]chat.Chat{},
		}
		svc := service.New(mockRepo)

		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		chats, err := svc.Get(ctx, "user2")
		if err != nil {
			t.Fatalf("failed to get chats: %v", err)
		}

		if len(chats) != 0 {
			t.Fatalf("expected 0 chats, got %d", len(chats))
		}
	})
}

func TestService_Chat(t *testing.T) {
	t.Run("should save chat successfully", func(t *testing.T) {
		mockRepo := &mockRepository{
			chats: map[string][]chat.Chat{},
		}
		svc := service.New(mockRepo)

		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		err := svc.Chat(ctx, "user2", "Hello, user2!")
		if err != nil {
			t.Fatalf("failed to save chat: %v", err)
		}

		if len(mockRepo.chats["user2"]) != 1 {
			t.Fatalf("expected 1 chat, got %d", len(mockRepo.chats["user2"]))
		}

		if mockRepo.chats["user2"][0].Message != "Hello, user2!" {
			t.Fatalf("expected message 'Hello, user2!', got '%s'", mockRepo.chats["user2"][0].Message)
		}
	})

	t.Run("should return error when repository Save fails", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}
		svc := service.New(mockRepo)

		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		err := svc.Chat(ctx, "user2", "Hello, user2!")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		expectedErr := "failed to save chat: repository error"
		if err.Error() != expectedErr {
			t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})
}
