package handler

import (
	"roadbarber/api/internal/models"
	"roadbarber/api/internal/modules/common/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	locationService *service.LocationService
}

func NewLocationHandler() *LocationHandler {
	return &LocationHandler{
		locationService: &service.LocationService{},
	}
}

// List 获取省份列表
func (h *LocationHandler) List(c *fiber.Ctx) error {
	locations, err := h.locationService.ListProvinces()
	if err != nil {
		return response.ServerError(c, "查询地区失败")
	}
	return response.Success(c, locations)
}

// GetChildren 获取下级地区
func (h *LocationHandler) GetChildren(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "地区ID不能为空")
	}

	locations, err := h.locationService.ListByParentID(id)
	if err != nil {
		return response.ServerError(c, "查询地区失败")
	}
	return response.Success(c, locations)
}

// Tree 地区树（用于首次加载完整数据）
func (h *LocationHandler) Tree(c *fiber.Ctx) error {
	locations, err := h.locationService.GetTree()
	if err != nil {
		return response.ServerError(c, "查询地区失败")
	}

	// 构建树形结构
	tree := buildTree(locations, "")
	return response.Success(c, tree)
}

// LocationNode 地区树节点
type LocationNode struct {
	models.Location
	Children []LocationNode `json:"children"`
}

func buildTree(locations []models.Location, parentID string) []LocationNode {
	var nodes []LocationNode
	for _, loc := range locations {
		var pid *string
		if loc.ParentID != nil {
			p := *loc.ParentID
			pid = &p
		}
		if (pid == nil && parentID == "") || (pid != nil && *pid == parentID) {
			children := buildTree(locations, loc.ID)
			nodes = append(nodes, LocationNode{
				Location: loc,
				Children: children,
			})
		}
	}
	return nodes
}
