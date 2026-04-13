package profile

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProfileRoutes(r fiber.Router, db *gorm.DB, jwtSecret string) {

	repo := NewRepo(db)
	service := NewProfileService(repo)
	controller := NewProfileController(service)

	profile := r.Group("profile")

	profile.Use(middleware.AuthMiddleWare(jwtSecret))
	profile.Get("/getprofile", controller.Profile)
	profile.Post("/editprofile", controller.EditProfile)
	profile.Post("/changepassword", controller.ChangePassWord)
	profile.Post("/deleteaccount", controller.DeleteAccount)

}