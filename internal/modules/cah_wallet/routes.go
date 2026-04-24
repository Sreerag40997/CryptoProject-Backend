package cashwallet

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router,service Service, jwtSecret string) {
	handler := NewHandler(service)

	wallet := app.Group("/wallet", middleware.AuthMiddleWare(jwtSecret))

	// user
	wallet.Post("/set-pin", handler.SetPin)
	wallet.Post("/change-pin", handler.ChangePin)

	wallet.Get("/me", handler.GetMyWallet)
	wallet.Get("/balance", handler.GetBalance)
	wallet.Get("/transactions", handler.GetTransactions)

	wallet.Post("/deposit", handler.Deposit)
	wallet.Post("/withdraw", handler.Withdraw)

	// admin
	admin := app.Group("/admin/wallet",middleware.AuthMiddleWare(jwtSecret), middleware.RequireRole("admin"))

	admin.Post("/:userId/block", handler.BlockWallet)
	admin.Post("/:userId/freeze", handler.FreezeWallet)
	admin.Post("/:userId/unblock", handler.UnblockWallet)

	admin.Post("/:userId/credit", handler.AdminCredit)

  app.Post("/webhook/razorpay", handler.RazorpayWebhook)
}
