package model

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	Trigger       string    `json:"trigger"`
	Routine       string    `json:"routine"`
	Reward        string    `json:"reward"`
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Completed     bool      `json:"completed"`
	StartDateTime time.Time `json:"start_date_time"`
	EndDateTime   time.Time `json:"end_date_time"`
}
