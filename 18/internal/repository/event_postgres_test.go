package repository

import (
	"testing"
	"time"
	"wbtech_l2/18/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	// Valid data
	id1, err1 := repo.Event.Create(1, model.Event{
		Description: "test_data",
		Date:        "2026-02-05",
		Time:        "23:00",
	})

	// Valid data (ID field is redundant)
	id2, err2 := repo.Event.Create(1, model.Event{
		ID:          1,
		Description: "test_data",
		Date:        "2026-02-05",
		Time:        "23:00",
	})

	// Valid data (Time is given with seconds)
	id3, err3 := repo.Event.Create(1, model.Event{
		Description: "test_data",
		Date:        "2026-02-05",
		Time:        "23:00",
	})

	// Invalid data (incorrect Date)
	id4, err4 := repo.Event.Create(1, model.Event{
		Description: "test_data",
		Date:        "date",
		Time:        "23:00",
	})

	// Invalid data (incorrect Time)
	id5, err5 := repo.Event.Create(1, model.Event{
		Description: "test_data",
		Date:        "2026-02-05",
		Time:        "time",
	})

	assert.NoError(t, err1)
	assert.NotEmpty(t, id1)

	assert.NoError(t, err2)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, 1, id2)

	assert.NoError(t, err3)
	assert.NotEmpty(t, id3)

	assert.Error(t, err4)
	assert.Empty(t, id4)

	assert.Error(t, err5)
	assert.Empty(t, id5)
}

func TestUpdateEvent(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	userID := 1

	id, _ := repo.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})

	// Valid
	err1 := repo.Event.Update(id, model.Event{
		Description: "my birthday",
		Date:        "2026-03-24",
		Time:        "16:00",
	})

	// Valid (updating just date)
	err2 := repo.Event.Update(id, model.Event{
		Date: "2026-03-24",
	})

	// Valid (updating just time)
	err3 := repo.Event.Update(id, model.Event{
		Time: "16:00",
	})

	// Valid (updating just description)
	err4 := repo.Event.Update(id, model.Event{
		Description: "my birthday",
	})

	// Invalid date
	err5 := repo.Event.Update(id, model.Event{
		Description: "my birthday",
		Date:        "date",
		Time:        "16:00",
	})

	// Invalid time
	err6 := repo.Event.Update(id, model.Event{
		Description: "my birthday",
		Date:        "2026-03-24",
		Time:        "time",
	})

	// Invalid eventID
	err7 := repo.Event.Update(12345, model.Event{
		Description: "my birthday",
		Date:        "2026-03-24",
		Time:        "16:00",
	})

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	assert.NoError(t, err4)

	assert.Error(t, err5)
	assert.Error(t, err6)
	assert.Error(t, err7)
	assert.Equal(t, NotFoundError, err7)
}

func TestDeleteEvent(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	userID := 1
	id1, _ := repo.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})
	id2, _ := repo.Event.Create(userID, model.Event{
		Description: "test_data",
		Date:        "2026-02-06",
		Time:        "14:00",
	})

	err1 := repo.Event.Delete(userID, id1)
	err2 := repo.Event.Delete(12345, id2)
	err3 := repo.Event.Delete(userID, 12345)

	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)
}

func TestGetEventsForDay(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	userID := 1
	eventDate := "2026-02-05"
	eventTime := "23:00"
	eventDescription := "test_data"
	eventsAmount := 3

	for i := 0; i < eventsAmount; i++ {
		repo.Event.Create(userID, model.Event{
			Description: eventDescription,
			Date:        eventDate,
			Time:        eventTime,
		})
	}

	eventsForDay1, err1 := repo.Event.GetEventsForDay(userID, eventDate)
	eventsForDay2, err2 := repo.Event.GetEventsForDay(userID, "2000-01-01")
	eventsForDay3, err3 := repo.Event.GetEventsForDay(12345, eventDate)

	assert.Equal(t, eventsAmount, len(eventsForDay1))

	assert.NoError(t, err1)
	assert.Equal(t, eventDate, eventsForDay1[0].Date)
	assert.Equal(t, eventDescription, eventsForDay1[0].Description)
	temp, _ := time.Parse("15:04", eventTime)
	assert.Equal(t, temp.Format("15:04:05"), eventsForDay1[0].Time)

	assert.NoError(t, err2)
	assert.Empty(t, eventsForDay2)

	assert.NoError(t, err3)
	assert.Empty(t, eventsForDay3)
}

func TestGetEventsForWeek(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	userID := 1
	eventDate := "2026-02-05"
	stillHaveEventDate := "2026-01-30"
	noEventDate := "2026-01-29"
	eventTime := "23:00"
	eventDescription := "test_data"
	eventsAmount := 3

	for i := 0; i < eventsAmount; i++ {
		repo.Event.Create(userID, model.Event{
			Description: eventDescription,
			Date:        eventDate,
			Time:        eventTime,
		})
	}

	eventsForWeek1, err1 := repo.Event.GetEventsForWeek(userID, eventDate)
	eventsForWeek2, err2 := repo.Event.GetEventsForWeek(userID, "2000-01-01")
	eventsForWeek3, err3 := repo.Event.GetEventsForWeek(12345, eventDate)
	eventsForWeek4, err4 := repo.Event.GetEventsForWeek(userID, noEventDate)
	eventsForWeek5, err5 := repo.Event.GetEventsForWeek(userID, stillHaveEventDate)

	assert.Equal(t, eventsAmount, len(eventsForWeek1))

	assert.NoError(t, err1)
	assert.Equal(t, eventDate, eventsForWeek1[0].Date)
	assert.Equal(t, eventDescription, eventsForWeek1[0].Description)
	temp, _ := time.Parse("15:04", eventTime)
	assert.Equal(t, temp.Format("15:04:05"), eventsForWeek1[0].Time)

	assert.NoError(t, err2)
	assert.Empty(t, eventsForWeek2)

	assert.NoError(t, err3)
	assert.Empty(t, eventsForWeek3)

	assert.NoError(t, err4)
	assert.Empty(t, eventsForWeek4)

	assert.NoError(t, err5)
	assert.NotEmpty(t, eventsForWeek5)
	assert.Equal(t, eventDate, eventsForWeek5[0].Date)
	assert.Equal(t, eventDescription, eventsForWeek5[0].Description)
	temp, _ = time.Parse("15:04", eventTime)
	assert.Equal(t, temp.Format("15:04:05"), eventsForWeek5[0].Time)
}

func TestGetEventsForMonth(t *testing.T) {
	db, teardown := TestDB(t)
	defer teardown(eventsTable)

	repo := NewRepository(db)

	userID := 1
	eventDate := "2026-02-05"
	stillHaveEventDate := "2026-01-06"
	noEventDate := "2026-01-05"
	eventTime := "23:00"
	eventDescription := "test_data"
	eventsAmount := 3

	for i := 0; i < eventsAmount; i++ {
		repo.Event.Create(userID, model.Event{
			Description: eventDescription,
			Date:        eventDate,
			Time:        eventTime,
		})
	}

	eventsForMonth1, err1 := repo.Event.GetEventsForMonth(userID, eventDate)
	eventsForMonth2, err2 := repo.Event.GetEventsForMonth(userID, "2000-01-01")
	eventsForMonth3, err3 := repo.Event.GetEventsForMonth(12345, eventDate)
	eventsForMonth4, err4 := repo.Event.GetEventsForMonth(userID, noEventDate)
	eventsForMonth5, err5 := repo.Event.GetEventsForMonth(userID, stillHaveEventDate)

	assert.Equal(t, eventsAmount, len(eventsForMonth1))

	assert.NoError(t, err1)
	assert.Equal(t, eventDate, eventsForMonth1[0].Date)
	assert.Equal(t, eventDescription, eventsForMonth1[0].Description)
	temp, _ := time.Parse("15:04", eventTime)
	assert.Equal(t, temp.Format("15:04:05"), eventsForMonth1[0].Time)

	assert.NoError(t, err2)
	assert.Empty(t, eventsForMonth2)

	assert.NoError(t, err3)
	assert.Empty(t, eventsForMonth3)

	assert.NoError(t, err4)
	assert.Empty(t, eventsForMonth4)

	assert.NoError(t, err5)
	assert.NotEmpty(t, eventsForMonth5)
	assert.Equal(t, eventDate, eventsForMonth5[0].Date)
	assert.Equal(t, eventDescription, eventsForMonth5[0].Description)
	temp, _ = time.Parse("15:04", eventTime)
	assert.Equal(t, temp.Format("15:04:05"), eventsForMonth5[0].Time)
}
