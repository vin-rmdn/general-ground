package repository

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/vin-rmdn/general-ground/internal/chat"
)

type repository struct {
	chats map[string][]chat.Chat
}

func New() repository {
	return repository{
		chats: make(map[string][]chat.Chat),
	}
}

func (r repository) Get(ctx context.Context, to string) ([]chat.Chat, error) {
	from, ok := ctx.Value(chat.FromKey{}).(string)
	if !ok {
		return nil, errors.New("fromUserID not found in context")
	}

	key := indexKey(from, to)
	messages, ok := r.chats[key]
	if !ok {
		return nil, errors.New("no chat found")
	}

	return messages, nil
}

func (r repository) Save(ctx context.Context, to, message string, timestamp time.Time) error {
	from, ok := ctx.Value(chat.FromKey{}).(string)
	if !ok {
		return errors.New("fromUserID not found in context")
	}

	key := indexKey(from, to)
	r.chats[key] = append(r.chats[key], chat.Chat{
		From:      from,
		To:        to,
		Message:   message,
		Timestamp: timestamp,
	})

	return nil
}

func indexKey(from, to string) string {
	users := []string{from, to}
	sort.Strings(users)

	return users[0] + "-" + users[1]
}
