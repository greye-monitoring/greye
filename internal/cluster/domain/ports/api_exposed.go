package ports

import "github.com/gofiber/fiber/v2"

type ApiExposed interface {
	Status(ctx *fiber.Ctx) error
	UpdateStatus(ctx *fiber.Ctx) error
	UpdateSingleStatus(ctx *fiber.Ctx) error
	Remove(ctx *fiber.Ctx) error
}
