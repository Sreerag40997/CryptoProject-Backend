package kyc

import (
	middleware "cryptox/internal/middleWare"
	cashwallet "cryptox/internal/modules/cah_wallet"
	ecard "cryptox/internal/modules/e_card"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring

	// wallet
	cashWalletRepo:=cashwallet.NewRepository(db)
	cashWalletService:=cashwallet.NewService(cashWalletRepo)

	// ecard 
	ecardRepo:=ecard.NewRepository(db)
	ecardService:=ecard.NewService(ecardRepo)

	// kyc
	repo := NewRepository(db)
	service := NewService(repo,cashWalletService,ecardService)
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

