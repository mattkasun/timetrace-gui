package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID      uuid.UUID
	Name    string
	Active  bool
	Updated time.Time
}
