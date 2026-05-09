// package centralroutes

// import (
// 	"cryptox/internal/modules/auth"
// 	cashwallet "cryptox/internal/modules/cah_wallet"
// 	cryptowallet "cryptox/internal/modules/crypto_wallet"
// 	ecard "cryptox/internal/modules/e_card"
// 	"cryptox/internal/modules/kyc"
// 	"cryptox/internal/modules/payment"
// 	tradeengine "cryptox/internal/modules/trade_engine"
// 	"cryptox/internal/modules/trade_engine/engine"

// 	"cryptox/internal/modules/profile"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/redis/go-redis/v9"
// 	"gorm.io/gorm"
// )

// func SetUp(app *fiber.App, db *gorm.DB, rdb *redis.Client, jwtSecret,razorpayKey,razorpaySecret string) {

// 	api := app.Group("/api")

// 	// payment
// 	paymentService := payment.NewRazorpayService(razorpayKey, razorpaySecret)

// 	// wallet
// 	walletRepo := cashwallet.NewRepository(db)
// 	walletService := cashwallet.NewService(walletRepo, paymentService)

// 	// E-card
// 	ecardRepo := ecard.NewRepository(db)
// 	ecardService := ecard.NewService(ecardRepo)

// 	// KYC 
// 	kycRepo := kyc.NewRepository(db)
// 	kycService := kyc.NewService(kycRepo, walletService, ecardService)

// 	// crypto wallet
// 	cryptoRepo:=cryptowallet.NewRepository(db)
// 	cryptoService:=cryptowallet.NewService(cryptoRepo)

// 	// trade engine
// 	traderepo:=tradeengine.NewRepository(db)
	
// 	enginExicutor:=engine.NewExecutor(traderepo)
// 	tradeService:=tradeengine.NewService(traderepo,rdb,enginExicutor)

// 	// routes
// 	auth.AuthRoutes(api, db, rdb, jwtSecret)
// 	profile.ProfileRoutes(api, db, jwtSecret)

// 	kyc.RegisterRoutes(api, kycService, jwtSecret)
// 	ecard.RegisterRoutes(api, ecardService, jwtSecret)
// 	cashwallet.RegisterRoutes(api, walletService, jwtSecret)
// 	cryptowallet.RegisterRoutes(api,cryptoService,jwtSecret)
// 	tradeengine.RegisterRoutes(api,tradeService,jwtSecret)
	
// }


package centralroutes

import (
	"cryptox/internal/modules/auth"
	cashwallet "cryptox/internal/modules/cah_wallet"
	cryptowallet "cryptox/internal/modules/crypto_wallet"
	ecard "cryptox/internal/modules/e_card"
	"cryptox/internal/modules/kyc"
	"cryptox/internal/modules/market"
	"cryptox/internal/modules/payment"
	"cryptox/internal/modules/profile"

	tradeengine "cryptox/internal/modules/trade_engine"
	"cryptox/internal/modules/trade_engine/bot"
	"cryptox/internal/modules/trade_engine/engine"

	walletadapter "cryptox/internal/modules/wallet_adapter"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(app *fiber.App, db *gorm.DB, rdb *redis.Client, jwtSecret, razorpayKey, razorpaySecret string) {

	api := app.Group("/api")

	// payment
	paymentService := payment.NewRazorpayService(razorpayKey, razorpaySecret)

	// cash wallet
	walletRepo := cashwallet.NewRepository(db)
	cashWalletService := cashwallet.NewService(walletRepo, paymentService)

	// crypto wallet
	cryptoRepo := cryptowallet.NewRepository(db)
	cryptoWalletService := cryptowallet.NewService(cryptoRepo)

	// wallet adapter (IMPORTANT)
	walletAdapter:= walletadapter.New(cashWalletService, cryptoWalletService)

	// ecard
	ecardRepo := ecard.NewRepository(db)
	ecardService := ecard.NewService(ecardRepo)

	// kyc
	kycRepo := kyc.NewRepository(db)
	kycService := kyc.NewService(kycRepo, cashWalletService, ecardService)

	// TRADE MODULE (CORRECT FLOW)

	tradeRepo := tradeengine.NewRepository(db)

	executor := engine.NewExecutor(tradeRepo, walletAdapter)

	eng := engine.NewEngine(executor)
	eng.Start()

	tradeService := tradeengine.NewService(tradeRepo, rdb, eng)

	// bot
	b := &bot.Bot{
		Engine: eng,
		Redis:  rdb,
	}
	b.Start()

	// routes
	auth.AuthRoutes(api, db, rdb, jwtSecret)
	profile.ProfileRoutes(api, db, jwtSecret)

	kyc.RegisterRoutes(api, kycService, jwtSecret)
	ecard.RegisterRoutes(api, ecardService, jwtSecret)
	cashwallet.RegisterRoutes(api, cashWalletService, jwtSecret)
	cryptowallet.RegisterRoutes(api, cryptoWalletService, jwtSecret)

	tradeengine.RegisterRoutes(api, tradeService, jwtSecret)
	market.RegisterRoutes(api,rdb)
}
