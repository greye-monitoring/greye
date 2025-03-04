package ports

import "github.com/gofiber/fiber/v2"

type ApiExposed interface {
	//AddApplicationBySvc(ctx *fiber.Ctx) error
	UnscheduleApplication(ctx *fiber.Ctx) error
	MonitoringApplication(ctx *fiber.Ctx) error
	GetApplicationMonitored(ctx *fiber.Ctx) error
	GetApplicationMonitoredByPod(ctx *fiber.Ctx) error
}
