package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/vin-rmdn/general-ground/internal/chat"
	"github.com/vin-rmdn/general-ground/internal/chat/repository"
)

func TestRepository_Get(t *testing.T) {
	t.Run("should error when no fromUserID in context", func(t *testing.T) {
		repo := repository.New()
		ctx := context.Background()
		chats, err := repo.Get(ctx, "user2")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		expectedErr := "fromUserID not found in context"
		if err.Error() != expectedErr {
			t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
		}
		if len(chats) != 0 {
			t.Fatalf("expected empty chats, got %v", chats)
		}
	})

	t.Run("should return empty chats when no chat exists", func(t *testing.T) {
		repo := repository.New()
		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")
		chats, err := repo.Get(ctx, "user2")

		if err == nil {
			t.Fatalf("error is not nil")
		}

		if len(chats) != 0 {
			t.Fatalf("expected empty chats, got %v", chats)
		}
	})

	t.Run("should return chats when chats exist", func(t *testing.T) {
		repo := repository.New()
		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")

		err := repo.Save(ctx, "user2", "Hello, user2!", time.Now())
		if err != nil {
			t.Fatalf("failed to save chat: %v", err)
		}

		chats, err := repo.Get(ctx, "user2")
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
}

func TestRepository_Save(t *testing.T) {
	t.Run("should save chat successfully", func(t *testing.T) {
		repo := repository.New()
		ctx := context.WithValue(context.Background(), chat.FromKey{}, "user1")

		err := repo.Save(ctx, "user2", "Hello, user2!", time.Now())
		if err != nil {
			t.Fatalf("failed to save chat: %v", err)
		}

		chats, err := repo.Get(ctx, "user2")
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

	t.Run("should return error when fromUserID is missing in context", func(t *testing.T) {
		repo := repository.New()
		err := repo.Save(context.Background(), "user2", "Hello, user2!", time.Now())

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		expectedErr := "fromUserID not found in context"
		if err.Error() != expectedErr {
			t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})
}
