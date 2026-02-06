package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wbtech_l2/18/internal/api/server"
	"wbtech_l2/18/internal/model"
	"wbtech_l2/18/internal/repository"
	"wbtech_l2/18/internal/service"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	testCases := []struct {
		name         string
		data         model.EventCreate
		expectedCode int
	}{
		{
			name: "valid",
			data: model.EventCreate{
				UserID:      1,
				Description: "something that I used to do",
				Date:        "2026-02-04",
				Time:        "14:55",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid date",
			data: model.EventCreate{
				UserID:      1,
				Description: "something that I used to do",
				Date:        "date",
				Time:        "14:55",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid (no date)",
			data: model.EventCreate{
				UserID:      1,
				Description: "something that I used to do",
				Time:        "14:55",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid time",
			data: model.EventCreate{
				UserID:      1,
				Description: "something that I used to do",
				Date:        "2026-02-06",
				Time:        "time",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid (no time)",
			data: model.EventCreate{
				UserID:      1,
				Description: "something that I used to do",
				Date:        "2026-02-06",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid (no description)",
			data: model.EventCreate{
				UserID: 1,
				Date:   "2026-02-06",
				Time:   "14:55",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid (no user_id)",
			data: model.EventCreate{
				Description: "something that I used to do",
				Date:        "2026-02-06",
				Time:        "14:55",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.data)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func TestUpdateEvent(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	userID := 1
	eventID, _ := repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})

	testCases := []struct {
		name         string
		data         model.Event
		expectedCode int
	}{
		{
			name: "valid",
			data: model.Event{
				ID:          eventID,
				Description: "something that I used to do",
				Date:        "2026-02-04",
				Time:        "14:55",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid date",
			data: model.Event{
				ID:          eventID,
				Description: "something that I used to do",
				Date:        "date",
				Time:        "14:55",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid time",
			data: model.Event{
				ID:          1,
				Description: "something that I used to do",
				Date:        "2026-02-06",
				Time:        "time",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid ID",
			data: model.Event{
				ID:          12345,
				Description: "something that I used to do",
				Date:        "2026-02-06",
				Time:        "14:55",
			},
			expectedCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.data)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/update_event", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func TestDeleteEvent(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	userID := 1
	eventID1, _ := repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})
	eventID2, _ := repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})

	testCases := []struct {
		name         string
		data         model.EventDelete
		expectedCode int
	}{
		{
			name: "valid",
			data: model.EventDelete{
				ID:     eventID1,
				UserID: userID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid (deleting same event twice)",
			data: model.EventDelete{
				ID:     eventID1,
				UserID: userID,
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name: "invalid (non-existing event_id)",
			data: model.EventDelete{
				ID:     12345,
				UserID: userID,
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name: "invalid (non-existing user_id)",
			data: model.EventDelete{
				ID:     eventID2,
				UserID: 12345,
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name: "invalid (no event id)",
			data: model.EventDelete{
				UserID: userID,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid (no user_id)",
			data: model.EventDelete{
				ID: eventID2,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.data)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/delete_event", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func TestGetEventsForDay(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	userID := 1
	eventDate := "2026-02-06"
	repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        eventDate,
		Time:        "14:00",
	})

	testCases := []struct {
		name         string
		params       string
		expectedCode int
	}{
		{
			name:         "valid",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (non-existing user_id)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", 12345, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (no events on given date)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, "2000-01-01"),
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid (no user_id)",
			params:       fmt.Sprintf("?date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid user_id",
			params:       fmt.Sprintf("?user_id=qwerty&date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid (no date)",
			params:       fmt.Sprintf("?user_id=%d", userID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid date",
			params:       fmt.Sprintf("?user_id=%d&date=qwerty", userID),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/events_for_day%s", tc.params), nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func TestGetEventsForWeek(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	userID := 1
	eventDate := "2026-02-06"
	repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        eventDate,
		Time:        "14:00",
	})

	testCases := []struct {
		name         string
		params       string
		expectedCode int
	}{
		{
			name:         "valid",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (non-existing user_id)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", 12345, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (no events on given date)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, "2000-01-01"),
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid (no user_id)",
			params:       fmt.Sprintf("?date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid user_id",
			params:       fmt.Sprintf("?user_id=qwerty&date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid (no date)",
			params:       fmt.Sprintf("?user_id=%d", userID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid date",
			params:       fmt.Sprintf("?user_id=%d&date=qwerty", userID),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/events_for_week%s", tc.params), nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func TestGetEventsForMonth(t *testing.T) {
	db, teardown := repository.TestDB(t)
	defer teardown()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)

	srv := new(server.Server)
	router := handlers.InitRoutes()
	go func() {
		if err := srv.Run("8888", router); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	userID := 1
	eventDate := "2026-02-06"
	repos.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        eventDate,
		Time:        "14:00",
	})

	testCases := []struct {
		name         string
		params       string
		expectedCode int
	}{
		{
			name:         "valid",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (non-existing user_id)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", 12345, eventDate),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid (no events on given date)",
			params:       fmt.Sprintf("?user_id=%d&date=%s", userID, "2000-01-01"),
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid (no user_id)",
			params:       fmt.Sprintf("?date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid user_id",
			params:       fmt.Sprintf("?user_id=qwerty&date=%s", eventDate),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid (no date)",
			params:       fmt.Sprintf("?user_id=%d", userID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid date",
			params:       fmt.Sprintf("?user_id=%d&date=qwerty", userID),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/events_for_month%s", tc.params), nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}
