package main

import (
	"FlexMadeTest/configuration"
	"FlexMadeTest/database"
	"FlexMadeTest/internal/repositories"
	"FlexMadeTest/internal/server/http"
	"FlexMadeTest/internal/services"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg, err := configuration.New()
	if err != nil {
		logger.WithField("error", err).Fatal("Error getting the App configuration!")
	}

	connection, err := database.NewGormDB(cfg)
	if err != nil {
		logger.WithField("error", err).Fatal("error DB connection!")
	}

	val := validator.New()
	r := repositories.NewStatistic(connection)
	s := services.NewStatistic(r)
	statisticTransport := http.NewStatistic(val, s)

	app := fiber.New()
	app.Get("/database/queries", statisticTransport.GetQueriesStatistic)

	if err = app.Listen(fmt.Sprintf(":%s", cfg.ApplicationPort)); err != nil {
		logger.WithField("error", err).Fatal("error starting http server")
	}
}
