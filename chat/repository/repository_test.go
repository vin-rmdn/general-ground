package repository_test

import (
	"context"
	"testing"

	"github.com/vin-rmdn/general-ground/chat"
	"github.com/vin-rmdn/general-ground/chat/repository"
)

func TestRepository_Get(t *testing.T) {
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
}