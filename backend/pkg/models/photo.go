package models

import "time"

type Photo struct {
	Path      string
	Tags      []string
	CreatedAt time.Time
	Data      string
}
