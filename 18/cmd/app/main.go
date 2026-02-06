package main

import (
	"wbtech_l2/18/internal/app"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Print("Initializing configs...")
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err.Error())
	}

	logrus.Print("Loading .env variables...")
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env variables: %s", err.Error())
	}

	app.Run()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
