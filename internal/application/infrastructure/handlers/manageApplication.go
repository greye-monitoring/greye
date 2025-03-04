package handlers

import (
	"github.com/gofiber/fiber/v2"
	model "greye/internal/application/domain/models"
	"net/url"
)

//// Register godoc
//// @Summary Register a new user
//// @Description Register a new user
//// @Accept json
//// @Produce json
//// @Tags Application
//// @Param body body model.RequestInfo true "User registration information"
//// @Success 200 {object} model.EasyResponse
//// @Failure 400 {object} model.EasyResponse
//// @Failure 500 {object} model.EasyResponse
//// @Router /api/v1/application [post]
//func (hdl *ApplicationHdl) AddApplicationBySvc(ctx *fiber.Ctx) error {
//	config, _ := hdl.config.GetConfig()
//
//	response := &model.EasyResponse{
//		Message: "Application added to monitoring successfully",
//	}
//	// retrieve hostname from env variable
//	hostname := os.Getenv("HOSTNAME")
//
//	//check if hostname finish with -0
//	r, _ := regexp.MatchString("-0$", hostname)
//	if !r {
//		hostname := fmt.Sprintf("%s-0.%s:%d", config.Server.ApplicationName, config.Server.ServiceHAName,
//			config.Server.Port)
//		//exec http call to hostname-0 to add application. put the same atom.Body
//
//		request := &models.HttpRequest{
//			Name:     hostname,
//			Host:     hostname,
//			Timeout:  5 * time.Second,
//			Protocol: "http",
//			Path:     "/api/v1/application",
//			Body:     ctx.Body(),
//			Method:   "POST",
//		}
//
//		response.Message = "Data sent to controller"
//
//		_, err := hdl.http.MakeRequest(request)
//		if err != nil {
//			hdl.logger.Error("Failed to make request to controller: ", err.Error())
//			response.Message = "Failed to make request to controller"
//
//		}
//		return ctx.Status(fiber.StatusOK).JSON(response)
//
//	}
//	var requestBody *v1.Service
//	if err := ctx.BodyParser(&requestBody); err != nil {
//		response.Message = "Invalid request body"
//		response.Error = err.Error()
//
//		return ctx.Status(fiber.StatusOK).JSON(response)
//	}
//	//data, _ := hdl.schedulerData.GetApplication("")
//	// Convert the map to JSON
//	//jsonData, err := json.Marshal(data)
//	//if err != nil {
//	//	hdl.logger.Error("Failed to marshal application data: ", err.Error())
//	//	response.Message = "Failed to marshal application data"
//	//	response.Error = err.Error()
//	//	return ctx.Status(fiber.StatusOK).JSON(response)
//	//}
//
//	//hdl.logger.Info(string(jsonData))
//	isToEnable := hdl.schedulerData.IsEnabled(requestBody)
//
//	// recupera l'oggett application dalle configurazioni generiche
//	defaultValue := config.Application
//
//	application := model.NewSchedulerApplicationFromService(requestBody, &defaultValue)
//	if err := hdl.validator.Struct(application); err != nil {
//		hdl.logger.Error(fmt.Sprintf("Error validating application: %v", err))
//		response.Message = "Error validating application"
//		response.Error = err.Error()
//		return ctx.Status(fiber.StatusOK).JSON(response)
//	}
//	if isToEnable {
//		err := hdl.schedulerData.AddApplication(application, false)
//		if err != nil {
//			response.Message = "error adding application to scheduler"
//			response.Error = err.Error()
//			return ctx.Status(fiber.StatusOK).JSON(response)
//		}
//	} else {
//		err := hdl.schedulerData.DeleteApplication(application)
//		if err != nil {
//			response.Message = "error deleting application from scheduler"
//			return ctx.Status(fiber.StatusOK).JSON(response)
//		}
//		response.Message = "Application deleted successfully"
//	}
//
//	return ctx.Status(fiber.StatusOK).JSON(response)
//}

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
