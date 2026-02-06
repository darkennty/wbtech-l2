package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"wbtech_l2/18/internal/model"

	"github.com/jmoiron/sqlx"
)

var NotFoundError = errors.New("event with given ID not found")

type EventPostgresRepository struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgresRepository {
	return &EventPostgresRepository{db: db}
}

func (r *EventPostgresRepository) Create(userID int, event model.Event) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, description, date, time) VALUES ($1, $2, $3, $4) RETURNING id;", eventsTable)
	row := r.db.QueryRow(query, userID, event.Description, event.Date, event.Time)
	err = row.Scan(&id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return 0, txErr
		}
		return 0, err
	}

	txErr := tx.Commit()
	if txErr != nil {
		return 0, txErr

	}
	return id, nil
}

func (r *EventPostgresRepository) Update(eventID int, event model.Event) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	fieldsToChange := make([]byte, 0)
	args := make([]interface{}, 0)

	if event.Description != "" {
		fieldsToChange = append(fieldsToChange, []byte(strings.Join([]string{fmt.Sprintf("description = $%d, ", len(args)+1)}, ""))...)
		args = append(args, event.Description)
	}

	if event.Date != "" {
		args = append(args, event.Date)
		fieldsToChange = append(fieldsToChange, []byte(strings.Join([]string{fmt.Sprintf("date = $%d, ", len(args))}, ""))...)
	}

	if event.Time != "" {
		args = append(args, event.Time)
		fieldsToChange = append(fieldsToChange, []byte(strings.Join([]string{fmt.Sprintf("time = $%d, ", len(args))}, ""))...)
	}

	toChangeStr := strings.TrimRight(string(fieldsToChange), ", ")
	args = append(args, eventID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d;", eventsTable, toChangeStr, len(args))
	affected, err := r.db.Exec(query, args...)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}

		return err
	}

	if temp, _ := affected.RowsAffected(); temp == 0 {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}

		return NotFoundError
	}

	txErr := tx.Commit()
	if txErr != nil {
		return txErr
	}

	return nil
}

func (r *EventPostgresRepository) Delete(userID, eventID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s e WHERE e.id = $1 AND e.user_id = $2;", eventsTable)
	affected, err := r.db.Exec(query, eventID, userID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}

		return err
	}

	if temp, _ := affected.RowsAffected(); temp == 0 {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}

		return NotFoundError
	}

	txErr := tx.Commit()
	if txErr != nil {
		return txErr
	}

	return nil
}

func (r *EventPostgresRepository) GetEventsForDay(userID int, date string) ([]model.Event, error) {
	var eventsFromDB []model.EventFromDB

	query := fmt.Sprintf("SELECT e.id, e.description, e.date, e.time FROM %s e WHERE e.user_id = $1 AND e.date = $2 ORDER BY e.time;", eventsTable)
	if err := r.db.Select(&eventsFromDB, query, userID, date); err != nil {
		return nil, err
	}

	events := make([]model.Event, 0, len(eventsFromDB))

	for _, dbEvent := range eventsFromDB {
		event := model.Event{
			ID:          dbEvent.ID,
			Description: dbEvent.Description,
			Date:        dbEvent.Date.Format("2006-01-02"),
			Time:        dbEvent.Time.Format("15:04:05"),
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventPostgresRepository) GetEventsForWeek(userID int, firstDate string) ([]model.Event, error) {
	var eventsFromDB []model.EventFromDB

	parsedDate, err := time.Parse("2006-01-02", firstDate)
	if err != nil {
		return nil, err
	}

	lastDate := parsedDate.Add(time.Hour * 24 * 7).Format("2006-01-02")

	query := fmt.Sprintf("SELECT e.id, e.description, e.date, e.time FROM %s e WHERE e.user_id = $1 AND e.date >= $2 AND e.date < $3 ORDER BY e.date, e.time;", eventsTable)
	if err = r.db.Select(&eventsFromDB, query, userID, firstDate, lastDate); err != nil {
		return nil, err
	}

	events := make([]model.Event, 0, len(eventsFromDB))

	for _, dbEvent := range eventsFromDB {
		event := model.Event{
			ID:          dbEvent.ID,
			Description: dbEvent.Description,
			Date:        dbEvent.Date.Format("2006-01-02"),
			Time:        dbEvent.Time.Format("15:04:05"),
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventPostgresRepository) GetEventsForMonth(userID int, firstDate string) ([]model.Event, error) {
	var eventsFromDB []model.EventFromDB

	parsedDate, err := time.Parse("2006-01-02", firstDate)
	if err != nil {
		return nil, err
	}

	lastDate := parsedDate.Add(time.Hour * 24 * 31).Format("2006-01-02")

	query := fmt.Sprintf("SELECT e.id, e.description, e.date, e.time FROM %s e WHERE e.user_id = $1 AND e.date >= $2 AND e.date < $3 ORDER BY e.date, e.time;", eventsTable)
	if err = r.db.Select(&eventsFromDB, query, userID, firstDate, lastDate); err != nil {
		return nil, err
	}

	events := make([]model.Event, 0, len(eventsFromDB))

	for _, dbEvent := range eventsFromDB {
		event := model.Event{
			ID:          dbEvent.ID,
			Description: dbEvent.Description,
			Date:        dbEvent.Date.Format("2006-01-02"),
			Time:        dbEvent.Time.Format("15:04:05"),
		}

		events = append(events, event)
	}

	return events, nil
}
