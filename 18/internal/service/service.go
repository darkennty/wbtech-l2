package service

import (
	"wbtech_l2/18/internal/model"
	"wbtech_l2/18/internal/repository"
)

type Event interface {
	Create(userID int, event model.Event) (int, error)
	Update(eventID int, event model.Event) error
	Delete(userID, eventID int) error
	GetEventsForDay(userID int, date string) ([]model.Event, error)
	GetEventsForWeek(userID int, date string) ([]model.Event, error)
	GetEventsForMonth(userID int, date string) ([]model.Event, error)
}

type Service struct {
	Event
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Event: NewEventService(repo.Event),
	}
}
