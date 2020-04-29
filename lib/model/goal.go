package model

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	User        string    `json:"user"`
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Achieved    bool      `json:"achieved"`
	Created     time.Time `json:"created"`
	Habits      []Habit   `json:"habits"`
}
