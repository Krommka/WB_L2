package usecase

import (
	"L2_18/internal/domain"
	"context"
)

type EventUC struct {
	repo Store
}

func NewEventUC(repo Store) *EventUC {
	return &EventUC{repo: repo}
}

func (uc *EventUC) CreateEvent(ctx context.Context, userUID string, event *domain.DTOEvent) error {

	user, err := uc.repo.GetUser(ctx, userUID)
	if err != nil {
		user, err = uc.repo.CreateUser(ctx)
	}
	if err != nil {
		return err
	}

	return uc.repo.CreateEvent(ctx, user.ID, event)
}

func (uc *EventUC) UpdateEvent(ctx context.Context, eventUID string, event *domain.DTOEvent) error {
	_, err := uc.repo.GetEventByID(ctx, eventUID)
	if err != nil {
		return err
	}
	err = uc.repo.UpdateEventByID(ctx, eventUID, event)

	return err
}

func (uc *EventUC) DeleteEventByID(ctx context.Context, eventUID string) error {
	return uc.repo.DeleteEventByID(ctx, eventUID)
}

func (uc *EventUC) GetEventByID(ctx context.Context, orderUID string) (*domain.Event, error) {
	return uc.repo.GetEventByID(ctx, orderUID)
}

func (uc *EventUC) GetEventForDay(ctx context.Context, userUID string, date string) ([]*domain.Event, error) {
	return uc.repo.GetEventForDay(ctx, userUID, date)
}

func (uc *EventUC) GetEventForWeek(ctx context.Context, userUID string, date string) ([]*domain.Event, error) {
	return uc.repo.GetEventForWeek(ctx, userUID, date)
}

func (uc *EventUC) GetEventForMonth(ctx context.Context, userUID string, date string) ([]*domain.Event, error) {
	return uc.repo.GetEventForMonth(ctx, userUID, date)
}
