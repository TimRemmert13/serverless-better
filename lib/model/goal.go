package model

import "time"

type Goal struct {
	User        string    `json:"user"`
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Achieved    bool      `json:"achieved"`
	Created     time.Time `json:"created"`
	Habits      []Habit   `json:"habits"`
}
