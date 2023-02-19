package entities

import "time"

type User struct {
	ID        int64
	VkID      int64
	CreatedAt time.Time
}
