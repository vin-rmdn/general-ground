package repository

import (
	"context"
	"errors"
	"sort"

	"github.com/vin-rmdn/general-ground/chat"
)

type repository struct {
	chats map[string][]chat.Chat
}

func New() repository {
	return repository{
		chats: make(map[string][]chat.Chat),
	}
}

func (r repository) Get(ctx context.Context, toUserID string) ([]chat.Chat, error) {
	fromUserID, ok := ctx.Value(chat.FromKey{}).(string)
	if !ok {
		return nil, errors.New("fromUserID not found in context")
	}

	users := []string{fromUserID, toUserID}
	sort.Strings(users)

	key := users[0] + "-" + users[1]
	messages, ok := r.chats[key]
	if !ok {
		return nil, errors.New("no chat found")
	}

	return messages, nil
}
