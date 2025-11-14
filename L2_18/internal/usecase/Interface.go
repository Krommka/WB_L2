package usecase

import (
	"L2_18/internal/domain"
	"context"
)

type Store interface {
	CreateUser(ctx context.Context) (domain.User, error)
	GetUser(ctx context.Context, userUID string) (domain.User, error)

	CreateEvent(ctx context.Context, userUID string, event *domain.DTOEvent) error
	UpdateEventByID(ctx context.Context, eventUID string, event *domain.DTOEvent) error
	DeleteEventByID(ctx context.Context, eventUID string) error
	GetEventByID(ctx context.Context, orderUID string) (domain.Event, error)
	GetEventForDay(ctx context.Context, userUID string, date string) ([]domain.Event, error)
	GetEventsForRange(ctx context.Context, userUID string, dateFrom string, dateTo string) ([]domain.Event, error)
}
