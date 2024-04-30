package api

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/danilotadeu/products/api/product"
	"github.com/danilotadeu/products/app"
	_ "github.com/danilotadeu/products/docs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

// validate is a validator package
var validate *validator.Validate

// @title		Star Wars API
// @version		1.0
// @BasePath	/api
func Register(apps *app.Container, port string) {
	fiberRoute := fiber.New()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	go func() {
		<-gracefulShutdown
		fmt.Println("Gracefully shutting down...")
		_ = fiberRoute.Shutdown()
	}()

	baseAPI := fiberRoute.Group("/api")

	validate = validator.New(validator.WithRequiredStructEnabled())

	// Planets
	product.NewAPI(baseAPI.Group("/products"), apps, validate)

	fiberRoute.Get("/swagger/*", swagger.HandlerDefault)

	logrus.WithFields(logrus.Fields{"trace": "api"}).Infof("Registered - Api")
	fiberRoute.Listen(":" + port)
}
