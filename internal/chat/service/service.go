package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/vin-rmdn/general-ground/internal/chat"
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

// Get is an exposed handler for getting messages.
// In my example, Get will be added with metrics and logs directly in the code
func (s service) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	from, _ := ctx.Value("user").(string)
	slog.Debug("get chat messages", slog.String("to", to), slog.String("from", from))

	chats, err := s.get(ctx, to)
	if err != nil {
		slog.Error(
			"failed to get chat messages",
			slog.String("to", to),
			slog.String("from", from),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("failed to get chat: %v", err)
	}

	slog.Debug(
		"successfully retrieved chat messages",
		slog.String("to", to),
		slog.String("from", from),
		slog.Any("chats", chats),
	)

	return chats, nil
}

// get is a private method that retrieves chat messages for a specific recipient.
// This method does not contain any logging or metrics and only contains business
// logic.
func (s service) get(ctx context.Context, to string) ([]chat.Chat, error) {
	// TODO: do something else here, like enriching user data

	chats, err := s.repository.Get(ctx, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %v", err)
	}

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
