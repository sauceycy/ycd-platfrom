package api

import (	"github.com/gofiber/fiber/v2"
	"ycd-platform/config"
)

// GetClusters 返回 Kubernetes 集群列表
func GetClusters(c *fiber.Ctx) error {
	clusters := config.Global.Clusters
	return c.JSON(clusters)
}

// GetNamespaces 返回指定集群的命名空间列表
func GetNamespaces(c *fiber.Ctx) error {
	clusterName := c.Query("cluster")
	if clusterName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing cluster query param",
		})
	}

	namespaces, err := fetchNamespacesForCluster(clusterName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(namespaces)
}

// 示例：获取命名空间的逻辑
func fetchNamespacesForCluster(clusterName string) ([]string, error) {
	// TODO: 实现与 Kubernetes API 的交互
	return []string{"default", "dev", "prod"}, nil
}