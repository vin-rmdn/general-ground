package chat

import "time"

type Chat struct {
	From      string
	To        string
	Message   string
	Timestamp time.Time
}
