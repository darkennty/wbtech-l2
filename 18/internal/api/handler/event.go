package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"wbtech_l2/18/internal/model"
	"wbtech_l2/18/internal/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createEvent(ctx *gin.Context) {
	var eventToCreate model.EventCreate
	if err := ctx.BindJSON(&eventToCreate); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "json contains incorrect data")
		return
	}

	if _, err := time.Parse("2006-01-02", eventToCreate.Date); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid date")
		return
	} else if _, err = time.Parse("15:04", eventToCreate.Time); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid time")
		return
	} else if eventToCreate.Description == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "no description given")
		return
	} else if eventToCreate.UserID == 0 {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "no user_id given")
		return
	}

	event := model.Event{
		Description: eventToCreate.Description,
		Date:        eventToCreate.Date,
		Time:        eventToCreate.Time,
	}

	id, err := h.services.Create(eventToCreate.UserID, event)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok", "id": id})
}

func (h *Handler) updateEvent(ctx *gin.Context) {
	var event model.Event
	if err := ctx.BindJSON(&event); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "json contains incorrect data")
		return
	}

	_, err := time.Parse("2006-01-02", event.Date)
	if event.Date != "" && err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid date")
		return
	}

	_, err = time.Parse("15:04", event.Time)
	if event.Time != "" && err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid time")
		return
	}

	err = h.services.Update(event.ID, event)
	if err != nil {
		if errors.Is(err, repository.NotFoundError) {
			ReturnErrorResponse(ctx, http.StatusServiceUnavailable, err.Error())
		} else {
			ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok"})
}

func (h *Handler) deleteEvent(ctx *gin.Context) {
	var eventDelete model.EventDelete
	if err := ctx.BindJSON(&eventDelete); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "json contains incorrect data")
		return
	}

	if eventDelete.UserID == 0 {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "no user_id given")
		return
	}

	if eventDelete.ID == 0 {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "no event id given")
		return
	}

	err := h.services.Delete(eventDelete.UserID, eventDelete.ID)
	if err != nil {
		if errors.Is(err, repository.NotFoundError) {
			ReturnErrorResponse(ctx, http.StatusServiceUnavailable, err.Error())
		} else {
			ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok"})
}

func (h *Handler) getEventsForDay(ctx *gin.Context) {
	stringUserID, ok := ctx.GetQuery("user_id")
	userID, err := strconv.Atoi(stringUserID)
	if !ok || stringUserID == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "user_id is required")
		return
	} else if err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, ok := ctx.GetQuery("date")
	if !ok || date == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "date is required")
		return
	} else if _, err = time.Parse("2006-01-02", date); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid date")
		return
	}

	var events []model.Event
	events, err = h.services.GetEventsForDay(userID, date)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok", "events": events})
}

func (h *Handler) getEventsForWeek(ctx *gin.Context) {
	stringUserID, ok := ctx.GetQuery("user_id")
	userID, err := strconv.Atoi(stringUserID)
	if !ok || stringUserID == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "user_id is required")
		return
	} else if err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, ok := ctx.GetQuery("date")
	if !ok || date == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "date is required")
		return
	} else if _, err = time.Parse("2006-01-02", date); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid date")
		return
	}

	var events []model.Event
	events, err = h.services.GetEventsForWeek(userID, date)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok", "events": events})
}

func (h *Handler) getEventsForMonth(ctx *gin.Context) {
	stringUserID, ok := ctx.GetQuery("user_id")
	userID, err := strconv.Atoi(stringUserID)
	if !ok || stringUserID == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "user_id is required")
		return
	} else if err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, ok := ctx.GetQuery("date")
	if !ok || date == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "date is required")
		return
	} else if _, err = time.Parse("2006-01-02", date); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid date")
		return
	}

	var events []model.Event
	events, err = h.services.GetEventsForMonth(userID, date)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ReturnResultResponse(ctx, gin.H{"status": "ok", "events": events})
}
