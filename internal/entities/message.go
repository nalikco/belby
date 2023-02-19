package entities

import "time"

type Message struct {
	ID        int64
	VkID      int64
	UserID    int64
	Message   string
	CreatedAt time.Time
}
