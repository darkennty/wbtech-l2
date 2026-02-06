package repository

import (
	"wbtech_l2/18/internal/model"

	"github.com/jmoiron/sqlx"
)

type Event interface {
	Create(userID int, event model.Event) (int, error)
	Update(eventID int, event model.Event) error
	Delete(userID, eventID int) error
	GetEventsForDay(userID int, date string) ([]model.Event, error)
	GetEventsForWeek(userID int, date string) ([]model.Event, error)
	GetEventsForMonth(userID int, date string) ([]model.Event, error)
}

type Repository struct {
	Event
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Event: NewEventPostgres(db),
	}
}
