package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	APIRouteBasePath = "/api/v1"

	HealthRoute = "/health"
)

type Controller interface {
	RegisterRoutes(router fiber.Router) error
}

func InitControllers() (controllers []Controller, err error) {
	controllers = append(controllers, NewHealthController())
	return controllers, nil
}

func SetupRoutes(v1Router fiber.Router, controllers []Controller) error {
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(v1Router); err != nil {
			return errors.Wrap(err, "failed to register routes")
		}
	}
	return nil
}
