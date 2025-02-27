package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	APIRouteBasePath = "/api/v1"

	TransactionRoute = "/transaction"
	AccountRoute     = "/account"

	HealthRoute = "/health"
)

type Controller interface {
	RegisterRoutes(router fiber.Router) error
}

func InitControllers() (controllers []Controller, err error) {
	controllers = append(controllers, NewHealthController())
	ledgerController, err := NewLedgerController()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init ledger controller")
	}
	controllers = append(controllers, ledgerController)
	return controllers, nil
}

func SetupRoutes(router fiber.Router, controllers []Controller) error {
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(router); err != nil {
			return errors.Wrap(err, "failed to register routes")
		}
	}
	return nil
}
