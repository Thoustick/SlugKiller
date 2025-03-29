package model

import "time"

type Link struct {
	ID        int64
	Slug      string
	URL       string
	CreatedAt time.Time
}
