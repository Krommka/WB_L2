package usecase

import (
	"L2_18/internal/domain"
	"context"
)

type store interface {
	CreateEvent(ctx context.Context, userID string, event *domain.Event) error
	UpdateEvent(ctx context.Context, order *domain.Event) (*domain.Event, error)
	DeleteEventByID(ctx context.Context, eventID string) error
	GetEventByID(ctx context.Context, orderUID string) (*domain.Event, error)
	GetEventForDay(ctx context.Context, userID string, date string) ([]*domain.Event, error)
	GetEventForWeek(ctx context.Context, userID string, date string) ([]*domain.Event, error)
	GetEventForMonth(ctx context.Context, userID string, date string) ([]*domain.Event, error)
}
