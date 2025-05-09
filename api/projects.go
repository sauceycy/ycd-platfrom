package api

import (
    "github.com/gofiber/fiber/v2"
    "ycd-platform/config"
)

// GetProjects 返回所有项目列表
func GetProjects(c *fiber.Ctx) error {
    projects := config.Global.Projects
    return c.JSON(projects)
}

// GetEnvironments 返回指定项目下的环境列表
func GetEnvironments(c *fiber.Ctx) error {
    projectName := c.Query("project")
    if projectName == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "missing project query param",
        })
    }

    for _, project := range config.Global.Projects {
        if project.Name == projectName {
            return c.JSON(project.Environments)
        }
    }

    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error": "project not found",
    })
}