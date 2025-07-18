package jobs

import (
	"time"
)

type Job struct {
	ID        int            `db:"id"`
	Type      string         `db:"type"`
	Payload   map[string]any `db:"payload"`
	Status    string         `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
