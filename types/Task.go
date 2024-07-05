package types

import (
	"time"
)

type Task struct {
	UserID      int       `json:"user-id"`
	ItemID      int       `json:"task-id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Bool        bool
}

func NewTask(userID, itemID int, desc string) *Task {
	return &Task{
		UserID:      userID,
		ItemID:      itemID,
		Description: desc,
		CreatedAt:   time.Now(),
	}
}
