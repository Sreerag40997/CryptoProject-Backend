package centralroutes

import (
	"cryptox/internal/modules/auth"
	cashwallet "cryptox/internal/modules/cah_wallet"
	cryptowallet "cryptox/internal/modules/crypto_wallet"
	ecard "cryptox/internal/modules/e_card"
	"cryptox/internal/modules/kyc"
	"cryptox/internal/modules/payment"

	"cryptox/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(app *fiber.App, db *gorm.DB, rdb *redis.Client, jwtSecret,razorpayKey,razorpaySecret string) {

	api := app.Group("/api")

	// PAYMENT (created once)
	paymentService := payment.NewRazorpayService(razorpayKey, razorpaySecret)

	// WALLET (depends on payment)
	walletRepo := cashwallet.NewRepository(db)
	walletService := cashwallet.NewService(walletRepo, paymentService)

	// ECARD
	ecardRepo := ecard.NewRepository(db)
	ecardService := ecard.NewService(ecardRepo)

	// KYC (depends on wallet + ecard)
	kycRepo := kyc.NewRepository(db)
	kycService := kyc.NewService(kycRepo, walletService, ecardService)

	// crypto wallet
	cryptoRepo:=cryptowallet.NewRepository(db)
	cryptoService:=cryptowallet.NewService(cryptoRepo)

	// routes
	auth.AuthRoutes(api, db, rdb, jwtSecret)
	profile.ProfileRoutes(api, db, jwtSecret)

	kyc.RegisterRoutes(api, kycService, jwtSecret)
	ecard.RegisterRoutes(api, ecardService, jwtSecret)
	cashwallet.RegisterRoutes(api, walletService, jwtSecret)
	cryptowallet.RegisterRoutes(api,cryptoService,jwtSecret)
	
}
