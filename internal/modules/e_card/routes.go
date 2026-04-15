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

	ecard.Get("/me",handler.GetMyCard)


	  // future features

	  // Block card
    // Regenerate card
    // Card transactions
    // Limits
    // Freeze/unfreeze

		// apis needed for future

		// GET    /card/me
    // POST   /card/block
    // POST   /card/unblock

}