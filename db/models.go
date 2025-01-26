package db

import "time"

type Todo struct {
	Ind         int
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}
