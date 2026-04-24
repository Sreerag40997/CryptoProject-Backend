package kyc

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, service Service, jwtSecret string) {
	//Dependency wiring

	handler := NewHandler(service)

	kyc:= app.Group("/kyc",middleware.AuthMiddleWare(jwtSecret))

  //user side
	kyc.Post("/submit",handler.SubmitKYC)
	kyc.Get("/status", handler.GetKYCStatus)
	kyc.Get("/me", handler.GetMyKYC)
  kyc.Put("/update", handler.UpdateKYC)

  // admin side

	admin:=app.Group("/admin/kyc",middleware.AuthMiddleWare(jwtSecret),middleware.RequireRole("admin"))

	admin.Get("/", handler.GetKYCList)
	admin.Get("/:id", handler.GetKYCByID)
	admin.Put("/:id", handler.UpdateKYCStatus)
}

