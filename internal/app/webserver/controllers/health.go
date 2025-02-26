package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) RegisterRoutes(router fiber.Router) error {
	router.Get(HealthRoute, func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	return nil
}
