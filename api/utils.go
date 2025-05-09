package api

import "github.com/gofiber/fiber/v2"

// ErrorResponse 通用错误响应
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
    return c.Status(status).JSON(fiber.Map{
        "error": message,
    })
}