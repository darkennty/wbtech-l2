package handler

import (
	"wbtech_l2/18/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/create_event", handlerFunc(h.createEvent))
	router.POST("/update_event", handlerFunc(h.updateEvent))
	router.POST("/delete_event", handlerFunc(h.deleteEvent))

	router.GET("/events_for_day", handlerFunc(h.getEventsForDay))
	router.GET("/events_for_week", handlerFunc(h.getEventsForWeek))
	router.GET("/events_for_month", handlerFunc(h.getEventsForMonth))

	return router
}
