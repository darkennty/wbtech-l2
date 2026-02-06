package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func handlerFunc(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Printf("%s %s %s\n", c.Request.Method, c.Request.RequestURI, time.Now().Format(time.RFC3339))
		f(c)
	}
}
