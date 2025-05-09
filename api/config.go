package api

import (
    "github.com/gofiber/fiber/v2"
    "ops-cd-platform/config"
)

// GetConfig 获取当前配置
func GetConfig(c *fiber.Ctx) error {
    return c.JSON(config.Global)
}

// UpdateProjectConfig 更新项目配置
func UpdateProjectConfig(c *fiber.Ctx) error {
    var updatedProject config.Project
    if err := c.BodyParser(&updatedProject); err != nil {
        return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
    }

    // 查找并更新项目
    for i, project := range config.Global.Projects {
        if project.Name == updatedProject.Name {
            config.Global.Projects[i] = updatedProject

            // 保存配置到文件
            if err := config.SaveConfig("config.yaml"); err != nil {
                return ErrorResponse(c, fiber.StatusInternalServerError, "failed to save config")
            }

            return c.JSON(fiber.Map{
                "message": "project updated successfully",
            })
        }
    }

    return ErrorResponse(c, fiber.StatusNotFound, "project not found")
}

// AddProjectConfig 添加新项目
func AddProjectConfig(c *fiber.Ctx) error {
    var newProject config.Project
    if err := c.BodyParser(&newProject); err != nil {
        return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
    }

    // 检查是否已存在
    for _, project := range config.Global.Projects {
        if project.Name == newProject.Name {
            return ErrorResponse(c, fiber.StatusConflict, "project already exists")
        }
    }

    // 添加新项目
    config.Global.Projects = append(config.Global.Projects, newProject)

    // 保存配置到文件
    if err := config.SaveConfig("config.yaml"); err != nil {
        return ErrorResponse(c, fiber.StatusInternalServerError, "failed to save config")
    }

    return c.JSON(fiber.Map{
        "message": "project added successfully",
    })
}

// DeleteProjectConfig 删除项目
func DeleteProjectConfig(c *fiber.Ctx) error {
    projectName := c.Query("name")
    if projectName == "" {
        return ErrorResponse(c, fiber.StatusBadRequest, "missing project name query param")
    }

    // 查找并删除项目
    for i, project := range config.Global.Projects {
        if project.Name == projectName {
            config.Global.Projects = append(config.Global.Projects[:i], config.Global.Projects[i+1:]...)

            // 保存配置到文件
            if err := config.SaveConfig("config.yaml"); err != nil {
                return ErrorResponse(c, fiber.StatusInternalServerError, "failed to save config")
            }

            return c.JSON(fiber.Map{
                "message": "project deleted successfully",
            })
        }
    }

    return ErrorResponse(c, fiber.StatusNotFound, "project not found")
}