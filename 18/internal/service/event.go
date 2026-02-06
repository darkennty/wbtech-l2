package service

import (
	"wbtech_l2/18/internal/model"
	"wbtech_l2/18/internal/repository"
)

type EventService struct {
	repo repository.Event
}

func NewEventService(repo repository.Event) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) Create(userID int, event model.Event) (int, error) {
	return s.repo.Create(userID, event)
}

func (s *EventService) Update(eventID int, event model.Event) error {
	return s.repo.Update(eventID, event)
}

func (s *EventService) Delete(userID, eventID int) error {
	return s.repo.Delete(userID, eventID)
}

func (s *EventService) GetEventsForDay(userID int, date string) ([]model.Event, error) {
	return s.repo.GetEventsForDay(userID, date)
}

func (s *EventService) GetEventsForWeek(userID int, date string) ([]model.Event, error) {
	return s.repo.GetEventsForWeek(userID, date)
}

func (s *EventService) GetEventsForMonth(userID int, date string) ([]model.Event, error) {
	return s.repo.GetEventsForMonth(userID, date)
}
