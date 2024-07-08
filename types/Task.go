package types

import (
	"time"
)

type Task struct {
	UserID      int       `json:"user-id"`
	ItemID      int       `json:"task-id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Date        time.Time `json:"date"`
}
