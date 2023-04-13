package models

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	ID      uuid.UUID
	Project string
	User    uuid.UUID
	Start   time.Time
	End     time.Time
}

func (r *Record) Duration() time.Duration {
	return r.End.Sub(r.Start)
}
