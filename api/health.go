package api

import "github.com/gofiber/fiber/v2"

// HealthCheck 健康检查接口
func HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
    })
}