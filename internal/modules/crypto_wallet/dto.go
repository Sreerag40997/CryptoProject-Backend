package cryptowallet

type CreateWalletRequest struct {
	Symbol string `json:"symbol"`
}

type AdminAmountRequest struct {
	Symbol string `json:"symbol"`
	Amount int64  `json:"amount"`
}

type CreateAssetRequest struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Precision int    `json:"precision"`
}