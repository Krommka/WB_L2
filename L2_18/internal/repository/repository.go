package repository

import (
	"L2_18/internal/domain"
	"context"
	"fmt"
	"github.com/jaevor/go-nanoid"
	"time"
)

const (
	retry = 5
)

type UserEventRepository struct {
	events map[string]*domain.Event
	users  map[string]*domain.User
}

func NewEventRepository() *UserEventRepository {
	events := make(map[string]*domain.Event)
	users := make(map[string]*domain.User)

	return &UserEventRepository{
		events: events,
		users:  users,
	}
}

func (repo *UserEventRepository) CreateEvent(ctx context.Context, userUID string,
	event *domain.DTOEvent) (string, error) {

	parsedDate, err := parseDate(event.Date)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	for range retry {
		generator, err := nanoid.Standard(12)
		if err != nil {
			continue
		}
		eventID := generator()
		if _, ok := repo.events[eventID]; ok {
			continue
		} else {
			repo.events[eventID] = &domain.Event{
				EventID:   eventID,
				CreatorID: userUID,
				Text:      event.Text,
				Date:      parsedDate,
			}
			return eventID, nil
		}
	}

	return "", fmt.Errorf("error creating new event. Retry is over")
}

func (repo *UserEventRepository) UpdateEvent(ctx context.Context, eventUID string,
	event *domain.DTOEvent) error {
	if _, ok := repo.events[eventUID]; !ok {
		return fmt.Errorf("event not found")
	}

	parsedDate, err := parseDate(event.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	repo.events[eventUID].Date = parsedDate
	repo.events[eventUID].Text = event.Text

	return nil

}
func (repo *UserEventRepository) DeleteEventByID(ctx context.Context, eventUID string) error {

	event, ok := repo.events[eventUID]
	if !ok {
		return fmt.Errorf("event not found")
	}
	userID := event.CreatorID
	if user, ok := repo.users[userID]; ok {
		delete(user.Events, eventUID)
	}
	delete(repo.events, eventUID)

	return nil
}
func (repo *UserEventRepository) GetEventByID(ctx context.Context, eventUID string) (domain.Event, error) {

	if _, ok := repo.events[eventUID]; !ok {
		return domain.Event{}, fmt.Errorf("event not found")
	}
	return *repo.events[eventUID], nil
}

func (repo *UserEventRepository) GetEventsForDay(ctx context.Context, userUID string, date string) ([]domain.Event,
	error) {
	parsedDate, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	user, ok := repo.users[userUID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	sl := make([]domain.Event, 0)
	for eventID := range user.Events {
		if event, ok := repo.events[eventID]; ok {
			if event.Date.Date() == parsedDate.Date() {
				sl = append(sl, *event)
			}
		}
	}

	return sl, nil
}

func (repo *UserEventRepository) GetEventsForRange(ctx context.Context, userUID string, dateFrom string,
	dateTo string) ([]domain.Event, error) {

	parsedDateFrom, err := parseDate(dateFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dateFrom: %v", err)
	}

	parsedDateTo, err := parseDate(dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dateTo: %v", err)
	}

	if parsedDateFrom.After(parsedDateTo) {
		return nil, fmt.Errorf("dateFrom is greater than dateTo")
	}

	user, ok := repo.users[userUID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	sl := make([]domain.Event, 0)

	for eventID := range user.Events {
		if event, ok := repo.events[eventID]; ok {
			if event.Date.After(parsedDateFrom) && event.Date.Before(parsedDateTo) {
				sl = append(sl, *event)
			}
		}
	}

	return sl, nil
}

func (repo *UserEventRepository) GetUser(ctx context.Context, userID string) (domain.User, error) {

	if _, ok := repo.users[userID]; !ok {
		return domain.User{}, fmt.Errorf("user not found")
	} else {
		return *repo.users[userID], nil
	}

}

func (repo *UserEventRepository) CreateUser(ctx context.Context) (string, error) {
	for range retry {
		generator, err := nanoid.Standard(8)
		if err != nil {
			continue
		}
		userID := generator()
		if _, ok := repo.users[userID]; ok {
			continue
		} else {
			user := &domain.User{
				ID:     userID,
				Events: make(map[string]struct{}),
			}
			repo.users[userID] = user
			return user.ID, nil
		}
	}
	return "", fmt.Errorf("error creating new user. Retry is over")
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
