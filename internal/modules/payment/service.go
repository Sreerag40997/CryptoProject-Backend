package payment

type Service interface {
	CreateOrder(amount int64, userID uint) (string, error)
	VerifySignature(orderID, paymentID, signature string) bool
	VerifyWebhookSignature(body []byte, signature string) bool
	FetchPaymentAmount(paymentID string) (int64, error)

	CreatePayout(userID uint, amount int64, name, ifsc, account string) (string, error)
}

