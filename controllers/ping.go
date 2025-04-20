package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type PingResponse struct {
	Message string `json:"message"`
}

// @Summary      Ping
// @Description  Responds with a "pong" message to test server availability
// @Tags         Health
// @Produce      json
// @Success      200  {object} controllers.PingResponse
// @Router       /v1/ping [get]
func HandlePing(c *fiber.Ctx) error {
	return c.JSON(PingResponse{
		Message: "pong",
	})
}
