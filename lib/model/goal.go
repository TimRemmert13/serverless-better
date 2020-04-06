package model

type Goal struct {
	User        string  `json:"user"`
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Title       string  `json:"title"`
	Achieved    bool    `json:"achieved"`
	Created     string  `json:"created"`
	Habits      []Habit `json:"habits"`
}
