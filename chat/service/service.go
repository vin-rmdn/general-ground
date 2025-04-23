package service

import (
	"context"
	"fmt"
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

func (s service) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	chats, err := s.repository.Get(ctx, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %v", err)
	}

	return chats, nil
}

func (s service) Chat(ctx context.Context, to, message string) error {
	if err := s.repository.Save(ctx, to, message, time.Now()); err != nil {
		return fmt.Errorf("failed to save chat: %v", err)
	}

	return nil
}
