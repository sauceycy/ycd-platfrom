package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// Deploy 触发 Helm 部署
func Deploy(c *fiber.Ctx) error {
    var payload struct {
        Project     string `json:"project"`
        Environment string `json:"environment"`
        Image       string `json:"image"`
        Cluster     string `json:"cluster"`
        Namespace   string `json:"namespace"`
    }
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    url := fmt.Sprintf("http://helm-api/%s/deploy", payload.Cluster)
    reqBody, _ := json.Marshal(payload)
    resp, err := http.Post(url, "application/json", strings.NewReader(string(reqBody)))
    if err != nil || resp.StatusCode != http.StatusOK {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to trigger deployment",
        })
    }

    return c.JSON(fiber.Map{
        "message": "deployment triggered successfully",
    })
}

// GetDeploymentStatus 查询部署状态
func GetDeploymentStatus(c *fiber.Ctx) error {
    cluster := c.Query("cluster")
    namespace := c.Query("namespace")
    release := c.Query("release")

    if cluster == "" || namespace == "" || release == "" {
        return ErrorResponse(c, fiber.StatusBadRequest, "missing required query parameters (cluster, namespace, release)")
    }

    // 调用 helm-api 查询部署状态
    url := fmt.Sprintf("http://helm-api/%s/status?namespace=%s&release=%s", cluster, namespace, release)
    resp, err := http.Get(url)
    if err != nil {
        return ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("failed to query deployment status: %v", err))
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return ErrorResponse(c, resp.StatusCode, "failed to get deployment status from helm-api")
    }

    // 返回 helm-api 的响应
    var status map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        return ErrorResponse(c, fiber.StatusInternalServerError, "failed to decode response from helm-api")
    }

    return c.JSON(status)
}