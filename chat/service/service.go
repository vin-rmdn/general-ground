package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/vin-rmdn/general-ground/chat"
)

type repository interface {
	Get(ctx context.Context, to string) ([]chat.Chat, error)
	Save(ctx context.Context, to, message string, timestamp time.Time) error
}

type service struct {
	repository repository
}

func New(repository repository) service {
	return service{
		repository: repository,
	}
}

// Get retrieves chat messages for a specific recipient.
// In my example, Get will be added with metrics directly in the code
func (s service) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	from, _ := ctx.Value("user").(string)
	slog.Debug("Get chat messages", slog.String("to", to), slog.String("from", from))

	chats, err := s.repository.Get(ctx, to)
	if err != nil {
		slog.Error(
			"Failed to get chat messages",
			slog.String("to", to),
			slog.String("from", from),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("failed to get chats: %v", err)
	}

	slog.Debug("Successfully retrieved chat messages", slog.String("to", to), slog.String("from", from), slog.Any("chats", chats))

	return chats, nil
}

// Chat sends a message to a specific recipient and saves it to the repository.
// In my example, Chat metrics will be added as a middleware
func (s service) Chat(ctx context.Context, to, message string) error {
	if err := s.repository.Save(ctx, to, message, time.Now()); err != nil {
		return fmt.Errorf("failed to save chat: %v", err)
	}

	return nil
}
