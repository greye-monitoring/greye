package handlers

import (
	"github.com/gofiber/fiber/v2"
	"greye/internal/cluster/domain/models"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Cluster
// @Success 200 {object} models.ClusterInfoDetails
// @Router /api/v1/cluster/status [get]
func (c ClusterHdl) Status(ctx *fiber.Ctx) error {
	status := c.cluster.ReadClustersStatuses()
	return ctx.Status(fiber.StatusOK).JSON(status)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Cluster
// @Param body body models.ClusterInfo true "User registration information"
// @Success 200 {object} models.ClusterInfo
// @Failure 400 {object} models.EasyResponses
// @Failure 500 {object} models.EasyResponses
// @Router /api/v1/cluster/status [put]
func (c ClusterHdl) UpdateStatus(ctx *fiber.Ctx) error {
	var requestBody models.ClusterInfoResponse

	//
	if err := ctx.BodyParser(&requestBody); err != nil {
		response := &models.EasyResponses{
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	status, err := c.cluster.Status(requestBody)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(status)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Cluster
// @Param body body models.SingleUpdateNode true "User registration information"
// @Success 200 {object} models.ClusterInfoDetails
// @Failure 400 {object} models.EasyResponses
// @Failure 500 {object} models.EasyResponses
// @Router /api/v1/cluster/suspend [put]
func (c ClusterHdl) UpdateSingleStatus(ctx *fiber.Ctx) error {
	var requestBody models.SingleUpdateNode

	if err := ctx.BodyParser(&requestBody); err != nil {
		response := &models.EasyResponses{
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	status, err := c.cluster.UpdateSingleNode(requestBody)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(status)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Tags Cluster
// @Success 200 {object} models.EasyResponses
// @Failure 500 {object} models.EasyResponses
// @Router /api/v1/cluster [delete]
func (c ClusterHdl) Remove(ctx *fiber.Ctx) error {

	isRemoved := c.cluster.Remove()
	if !isRemoved {
		response := &models.EasyResponses{
			Message: "Error removing cluster",
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := &models.EasyResponses{
		Message: "Cluster removed succesfully",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
