package models

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	ID      uuid.UUID
	Project uuid.UUID
	Start   time.Time
	End     time.Time
}

func (r *Record) Duration() time.Duration {
	return r.End.Sub(r.Start)
}
