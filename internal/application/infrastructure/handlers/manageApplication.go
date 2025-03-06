package handlers

import (
	"github.com/gofiber/fiber/v2"
	model "greye/internal/application/domain/models"
	"net/url"
)

// UnscheduleApplication è la funzione chiamata per gestire l'endpoint DELETE.
// @Summary Unschedule application monitoring
// @Description Rimuove il monitoraggio per un'applicazione specifica identificata dal suo nome.
// @Tags Application
// @Accept json
// @Produce json
// @Param service path string true "Nome del servizio da rimuovere"
// @Success 200 {object} model.EasyResponse
// @Failure 400 {object} model.EasyResponse
// @Failure 500 {object} model.EasyResponse
// @Router /api/v1/application/monitor/{service} [delete]
func (hdl *ApplicationHdl) UnscheduleApplication(ctx *fiber.Ctx) error {
	response := &model.EasyResponse{
		Message: "Service unscheduled",
	}

	service, err := url.QueryUnescape(ctx.Params("service"))
	if err != nil {
		response.Message = "Error decoding the service"
		response.Error = err.Error()
		hdl.logger.Error("%s", response.Message)
		hdl.logger.Error("%s", response.Error)
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = hdl.schedulerData.DeleteApplicationFromUrl(service)
	if err != nil {
		response.Message = "Error unscheduling application"
		response.Error = err.Error()
		hdl.logger.Error("%s", response.Message)
		hdl.logger.Error("%s", response.Error)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Application
// @Param url query string false "url" example(http://example.com)
// @Success 200 {object} map[string]models.SchedulerApplication
// @Failure 500 {object} model.EasyResponse
// @Router /api/v1/application/monitor [get]
func (hdl *ApplicationHdl) GetApplicationMonitored(ctx *fiber.Ctx) error {

	url := ctx.Query("url")
	data, err := hdl.schedulerData.GetApplication(url)
	if err != nil {
		hdl.logger.Error("Failed to get application data: ")
		response := &model.EasyResponse{
			Message: "Failed to retrieve application data",
			Error:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Application
// @Param body body model.SchedulerApplication true "User registration information"
// @Success 200 {object} model.EasyResponse
// @Failure 400 {object} model.EasyResponse
// @Failure 500 {object} model.EasyResponse
// @Router /api/v1/application/monitor [put]
func (hdl *ApplicationHdl) MonitoringApplication(ctx *fiber.Ctx) error {
	response := &model.EasyResponse{
		Message: "Application added to monitoring successfully",
	}

	var requestBody []model.SchedulerApplication
	if err := ctx.BodyParser(&requestBody); err != nil {
		response.Message = "Invalid request body"
		response.Error = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	for _, app := range requestBody {
		if err := hdl.validator.Struct(&app); err != nil {
			response.Message = "Error validating application"
			response.Error = err.Error()
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}

		application := model.GenerateJobSchedulerApplication(app)
		err := hdl.schedulerData.MonitorApplication(application, false)
		if err != nil {
			response.Message = "Error while adding application to monitor"
			response.Error = err.Error()
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Register godoc
// @Summary get monitored applications by pod
// @Description get all applications distinct by pod
// @Accept json
// @Produce json
// @Tags Application
// @Success 200 {object} map[string][]model.SchedulerApplication
// @Failure 500 {object} model.EasyResponse
// @Router /api/v1/application/monitor/pod [get]
func (hdl *ApplicationHdl) GetApplicationMonitoredByPod(ctx *fiber.Ctx) error {
	data, err := hdl.schedulerData.GetApplication("")
	if err != nil {
		hdl.logger.Error("Failed to get application data: ")
		response := &model.EasyResponse{
			Message: "Failed to retrieve application data",
			Error:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	mapData := make(map[string][]model.SchedulerApplication)
	for _, app := range data {
		mapData[app.ScheduledApplication] = append(mapData[app.ScheduledApplication], app)
	}
	return ctx.Status(fiber.StatusOK).JSON(mapData)
}
