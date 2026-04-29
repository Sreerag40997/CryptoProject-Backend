package cryptowallet

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, service Service, jwtSecret string) {

	h := NewHandler(service)

	// USER ROUTES
	wallet := app.Group("/crypto-wallet", middleware.AuthMiddleWare(jwtSecret))

	wallet.Post("/", h.CreateWallet)
	wallet.Get("/", h.GetWallets)
	wallet.Get("/summary", h.GetSummary)
	wallet.Get("/transactions", h.GetTransactions)
	wallet.Get("/locks", h.GetLocks)
	wallet.Get("/:symbol", h.GetWallet)

	// ADMIN ROUTES
	admin := app.Group("/admin/crypto-wallet",
		middleware.AuthMiddleWare(jwtSecret),
		middleware.RequireRole("admin"),
	)

	admin.Get("/", h.GetAllWalletsAdmin)
	admin.Get("/:userId", h.GetUserWalletsAdmin)
	admin.Get("/:userId/transactions", h.GetUserTransactionsAdmin)
	admin.Get("/:userId/:symbol", h.GetUserWalletBySymbolAdmin)

	admin.Post("/:userId/credit", h.AdminCredit)
	admin.Post("/:userId/debit", h.AdminDebit)

	admin.Post("/:userId/freeze", h.FreezeWallet)
	admin.Post("/:userId/unfreeze", h.UnfreezeWallet)

	admin.Get("/transactions", h.GetAllTransactionsAdmin)
	

	// ASSET ROUTES
	assetAdmin := app.Group("/admin/assets",
		middleware.AuthMiddleWare(jwtSecret),
		middleware.RequireRole("admin"),
	)

	assetAdmin.Post("/", h.CreateAsset)
	assetAdmin.Put("/:id", h.UpdateAsset)
	assetAdmin.Put("/:id/status", h.UpdateAssetStatus)

	app.Get("/assets", h.GetAssets)
}


