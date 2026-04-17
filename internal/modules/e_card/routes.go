package ecard

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	ecard := app.Group("/ecard",middleware.AuthMiddleWare(jwtSecret))

  //user side apis
	ecard.Get("/me",handler.GetMyCard)
	ecard.Post("/block",handler.BlockCard)
	ecard.Post("/unblock",handler.UnblockCard)

	//admin side apis

	admin:=app.Group("/admin/ecard",middleware.AuthMiddleWare(jwtSecret),middleware.RequireRole("admin"))

	admin.Get("/:userId", handler.AdminGetCard)
	admin.Post("/:userId/block", handler.AdminBlockCard)
	admin.Post("/:userId/unblock", handler.AdminUnblockCard)


}