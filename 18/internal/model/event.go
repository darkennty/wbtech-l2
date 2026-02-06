package model

import "time"

type Event struct {
	ID          int    `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`
	Time        string `json:"time" db:"time"`
}

type EventFromDB struct {
	ID          int       `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	Time        time.Time `json:"time" db:"time"`
}

type EventCreate struct {
	UserID      int    `json:"user_id" db:"user_id"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`
	Time        string `json:"time" db:"time"`
}

type EventDelete struct {
	ID     int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
}
