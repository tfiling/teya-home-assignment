package main

import (
	"fmt"
	"teya_home_assignment/internal/app/webserver/controllers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	apiGroup := app.Group(controllers.APIRouteBasePath)
	APIControllers, err := controllers.InitControllers()
	if err != nil {
		panic(fmt.Errorf("error setting up controllers: %w", err))
	}
	if err := controllers.SetupRoutes(apiGroup, APIControllers); err != nil {
		panic(fmt.Errorf("error setting up routes: %w", err))
	}

	if err := app.Listen(":8000"); err != nil {
		panic(fmt.Errorf("error starting server: %w", err))
	}
}
