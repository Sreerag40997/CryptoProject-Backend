package rbac

import (
	"cryptox/packages/utils"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// POST /api/admin/permissions
func (h *Controller) CreatePermission(c *fiber.Ctx) error {
	var req PermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid body", err)
	}

	if err := h.svc.AddNewPermission(req); err != nil {
		return utils.Error(c, 500, "", err.Error())
	}
	return utils.Success(c, 201, "Permission created", nil)
}

// PATCH /api/admin/roles/toggle
func (h *Controller) ToggleRolePermission(c *fiber.Ctx) error {
	var req AssignPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid body", err.Error())
	}

	if err := h.svc.TogglePermission(req); err != nil {
		return utils.Error(c, 500, "", err.Error())
	}
	return utils.Success(c, 200, "Role updated successfully", nil)
}

func (h *Controller) CreateRole(c *fiber.Ctx) error {

	var req CreateRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid body", err.Error())
	}

	if err := h.svc.CreateRole(req); err != nil {
		return utils.Error(c, 500, "", err.Error())
	}

	return utils.Error(c, 201, "role created", nil)
}

func (h *Controller) GetPermissionsByRole(c *fiber.Ctx) error {

	role := c.Params("role")

	data, err := h.svc.GetPermissionsByRole(role)
	if err != nil {
		return utils.Error(c, 500, "", err.Error())
	}

	return utils.Success(c, 200, "", data)
}

func (h *Controller) ListRoles(c *fiber.Ctx) error {
	roles, err := h.svc.GetAllRoles()
	if err != nil {
		return utils.Error(c, 500, "", err.Error())
	}
	return utils.Success(c, 200, "", roles)
}

func (h *Controller) ListPermissions(c *fiber.Ctx) error {
	permissions, err := h.svc.GetAllPermissions()
	if err != nil {
		return utils.Error(c, 500, "", err.Error())
	}
	return utils.Success(c, 200, "", permissions)
}