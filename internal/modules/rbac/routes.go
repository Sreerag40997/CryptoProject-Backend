package rbac

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	permission := r.Group("/admin", middleware.AuthMiddleWare(jwtSecret), middleware.RequireRole("admin"))

	permission.Post("/permission", controller.CreatePermission)
	permission.Get("/permissions", controller.ListPermissions)
	permission.Patch("/roles/toggle", controller.ToggleRolePermission)
	permission.Post("/role", controller.CreateRole)
	permission.Get("/roles", controller.ListRoles)
	permission.Get("/roles/:role/permissions", controller.GetPermissionsByRole)

}