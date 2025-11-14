package domain

import "time"

type Event struct {
	CreatorID string
	EventID   string
	Text      string
	Date      time.Time
}

type DTOEvent struct {
	Text string
	Date string
}
