package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/gofiber/fiber/v2"
    "ycd-platform/config"
)

// GetImages 实时从 Harbor 获取最近 10 个镜像 tag
func GetImages(c *fiber.Ctx) error {
    projectName := c.Query("project")
    if projectName == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "missing project query param",
        })
    }

    // 找到 harbor repo
    var harborRepo string
    for _, p := range config.Global.Projects {
        if p.Name == projectName {
            harborRepo = p.HarborRepo
            break
        }
    }
    if harborRepo == "" {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "project not found",
        })
    }

    // 解析 repo 路径
    split := strings.SplitN(harborRepo, "/", 2)
    if len(split) != 2 {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "invalid harbor repo format",
        })
    }
    project := split[0]
    repo := split[1]

    // 获取镜像 tags
    tags, err := fetchHarborTags(project, repo)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(tags)
}

// 请求 Harbor API 获取镜像 tag
func fetchHarborTags(project, repo string) ([]string, error) {
    url := fmt.Sprintf("https://hub.dream22.xyz/api/v2.0/projects/%s/repositories/%s/artifacts?page_size=50", project, repo)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // 从配置中读取认证信息
    username := config.Global.HarborAuth.Username
    password := config.Global.HarborAuth.Password
    req.SetBasicAuth(username, password)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to query Harbor, status code: %d", resp.StatusCode)
    }

    var data []struct {
        Tags []struct {
            Name string `json:"name"`
        } `json:"tags"`
    }
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    var tags []string
    for _, artifact := range data {
        for _, tag := range artifact.Tags {
            tags = append(tags, tag.Name)
        }
    }

    if len(tags) > 10 {
        tags = tags[:10]
    }

    return tags, nil
}