package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wbtech_l2/18/internal/api/handler"
	"wbtech_l2/18/internal/api/server"
	"wbtech_l2/18/internal/repository"
	"wbtech_l2/18/internal/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	logrus.Print("Initializing DB...")
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASS"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  viper.GetString("db.ssl_mode"),
	})
	if err != nil {
		logrus.Fatalf("Error initializing DB: %s", err.Error())
	}

	logrus.Print("Initializing components...")
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Error occured while running http-server: %s", err.Error())
		}
	}()

	logrus.Print("App started.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App is shutting down.")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Error occured while shutting down server: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Fatalf("Error occured while closing DB: %s", err.Error())
	}

	logrus.Print("App is stopped.")
}
